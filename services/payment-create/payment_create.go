package paycreate

import (
	"fmt"
	"time"

	"github.com/MihaiBlebea/purpletree/platform/discount"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	d "github.com/MihaiBlebea/purpletree/platform/discount"
	e "github.com/MihaiBlebea/purpletree/platform/email"
	"github.com/MihaiBlebea/purpletree/platform/payment"
	p "github.com/MihaiBlebea/purpletree/platform/product"
	u "github.com/MihaiBlebea/purpletree/platform/user"
	t "github.com/MihaiBlebea/purpletree/platform/user/token"
)

// New returns a CreatePaymentService struct
func New() *CreatePaymentService {
	conn := c.Mysql()

	userRepository := *u.Repo(conn)
	productRepository := *p.Repo(conn)
	tokenRepository := *t.Repo(conn)
	paymentRepository := *payment.Repo(conn)
	discountRepository := *d.Repo(conn)
	emailService := e.Service{}

	return &CreatePaymentService{
		userRepository,
		productRepository,
		tokenRepository,
		paymentRepository,
		discountRepository,
		emailService,
	}
}

// CreatePaymentService creates a payment
type CreatePaymentService struct {
	UserRepository     u.Repository
	ProductRepository  p.Repository
	TokenRepository    t.Repository
	PaymentRepository  payment.Repository
	DiscountRepository discount.Repository
	EmailService       e.Service
}

// CreatePaymentResponse is the response for CreatePaymentService
type CreatePaymentResponse struct {
	PaymentID   int     `json:"payment_id"`
	UserID      int     `json:"user_id"`
	ProductID   int     `json:"product_id"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ProductName string  `json:"product_name"`
	UserName    string  `json:"user_name"`
	Token       string  `json:"token"`
	Success     bool    `json:"success"`
}

// CreatePaymentRequest request DTO to create a payment
type CreatePaymentRequest struct {
	Token        string
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Nonce        string
	PaymentType  string
	ProductCode  string
	DiscountCode string
}

// Execute creates a payment when the user is not logged in the platform
func (s *CreatePaymentService) Execute(request CreatePaymentRequest) (response CreatePaymentResponse, err error) {
	// Figure out if the request implies creating a new user or user is already known
	// Does request has an auth token?
	user, err := s.getUser(request)
	if err != nil {
		return response, err
	}

	// Get product information from db
	product, _, err := s.ProductRepository.FindByCode(request.ProductCode)
	if err != nil {
		return response, err
	}

	// Get the discount if any
	discount, count, err := s.DiscountRepository.FindByCode(request.DiscountCode)
	if err != nil {
		return response, err
	}

	// Calculate the price and apply discount if any
	price := payment.NewPrice(product.Price, product.Currency)
	if count != 0 && discount.IsValid() == true {
		price.ApplyDiscount(discount.Percentage)
	}

	price.ApplyDiscount(discount.Percentage)

	// Make payment with payment provider
	paymentService := payment.Service{}
	params := payment.Params{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     *price,
		Nonce:     request.Nonce,
	}
	pay, err := paymentService.Pay(
		s.PaymentRepository,
		params,
	)
	if err != nil {
		return response, err
	}

	// Generate 3 activation tokens in user account
	tokens, err := t.NewCount(user.ID, user.Email, 3)
	if err != nil {
		return response, err
	}
	for _, token := range *tokens {
		_, err = s.TokenRepository.Add(&token)
		if err != nil {
			return response, err
		}
	}

	// Send an invoice email to the user
	err = s.EmailService.Send("payment", s.createEmailPayload(
		user.Name,
		user.Email,
		product.Name,
		pay.Created,
		price.GetAmount(),
		price.WithTVA(),
	))
	if err != nil {
		return response, err
	}

	return CreatePaymentResponse{
		PaymentID:   pay.ID,
		UserID:      user.ID,
		ProductID:   product.ID,
		Price:       price.WithTVA(),
		Currency:    product.Currency,
		ProductName: product.Name,
		UserName:    user.Name,
		Token:       user.JWT,
		Success:     true,
	}, nil
}

func (s *CreatePaymentService) createUserFromRaw(firstName, lastName, email, password string) (user *u.User, err error) {
	user, err = u.New(
		fmt.Sprintf("%s %s", firstName, lastName),
		email,
		password,
		true,
	)
	if err != nil {
		return user, err
	}

	// Add user to db
	userID, err := s.UserRepository.Add(user)
	if err != nil {
		return user, err
	}
	user.ID = userID

	return user, nil
}

func (s *CreatePaymentService) createEmailPayload(name, email, productName string, paymentDate time.Time, price, totalPrice float64) map[string]interface{} {
	payload := make(map[string]interface{})
	payload["name"] = name
	payload["email"] = email
	payload["productName"] = productName
	payload["paymentDate"] = paymentDate
	payload["price"] = price
	payload["totalPrice"] = totalPrice

	return payload
}

func (s *CreatePaymentService) getUser(request CreatePaymentRequest) (user *u.User, err error) {
	if request.hasToken() == true {
		user, count, err := s.UserRepository.FindByJWT(request.Token)
		if count == 0 || err != nil {
			return user, fmt.Errorf("Could not find user with token %s", request.Token)
		}

		return user, nil
	}

	user, err = s.createUserFromRaw(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *CreatePaymentRequest) hasToken() bool {
	if r.Token == "" {
		return false
	}

	return true
}
