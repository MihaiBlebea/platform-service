package paycreate

import (
	"fmt"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	e "github.com/MihaiBlebea/Wordpress/platform/email"
	"github.com/MihaiBlebea/Wordpress/platform/payment"
	p "github.com/MihaiBlebea/Wordpress/platform/product"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
	t "github.com/MihaiBlebea/Wordpress/platform/user/token"
)

// New returns a CreatePaymentService struct
func New() *CreatePaymentService {
	conn := c.Mysql()

	userRepository := *u.Repo(conn)
	productRepository := *p.Repo(conn)
	tokenRepository := *t.Repo(conn)
	paymentRepository := *payment.Repo(conn)
	emailService := e.Service{}

	return &CreatePaymentService{
		userRepository,
		productRepository,
		tokenRepository,
		paymentRepository,
		emailService,
	}
}

// CreatePaymentService creates a payment
type CreatePaymentService struct {
	UserRepository    u.Repository
	ProductRepository p.Repository
	TokenRepository   t.Repository
	PaymentRepository payment.Repository
	EmailService      e.Service
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

// ExecuteWithAuth creates a payment when the user is logged in the platform
func (s *CreatePaymentService) ExecuteWithAuth(token, nonce, paymentType, productCode string) (response CreatePaymentResponse, err error) {
	user, count, _ := s.UserRepository.FindByJWT(token)
	if count == 0 {
		return response, fmt.Errorf("Could not find user with token %s", token)
	}

	// Get product information from db
	product, _, err := s.ProductRepository.FindByCode(productCode)
	if err != nil {
		return response, err
	}

	// Make payment with payment provider
	amount := payment.NewAmount(product.Price, product.Currency)

	paymentService := payment.Service{}
	params := payment.Params{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     *amount,
		Nonce:     nonce,
	}
	pay, err := paymentService.Pay(
		s.PaymentRepository,
		params,
	)
	if err != nil {
		return response, err
	}

	// Generate activation tokens in user account
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

	// Send an invoice email
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["email"] = user.Email
	data["product"] = product.Name
	err = s.EmailService.Send("payment", data)
	if err != nil {
		return response, err
	}

	return CreatePaymentResponse{
		PaymentID:   pay.ID,
		UserID:      user.ID,
		ProductID:   product.ID,
		Price:       amount.GetFloat(),
		Currency:    product.Currency,
		ProductName: product.Name,
		UserName:    user.Name,
		Token:       user.JWT,
		Success:     true,
	}, nil
}

// Execute creates a payment when the user is not logged in the platform
func (s *CreatePaymentService) Execute(firstName, lastName, email, password, nonce, paymentType, productCode string) (response CreatePaymentResponse, err error) {
	user, err := u.New(
		fmt.Sprintf("%s %s", firstName, lastName),
		email,
		password,
		true,
	)
	if err != nil {
		return response, err
	}

	// Add user to db
	userID, err := s.UserRepository.Add(user)
	if err != nil {
		return response, err
	}
	user.ID = userID

	// Get product information from db
	product, _, err := s.ProductRepository.FindByCode(productCode)
	if err != nil {
		return response, err
	}

	// Make payment with payment provider
	amount := payment.NewAmount(product.Price, product.Currency)

	paymentService := payment.Service{}
	params := payment.Params{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     *amount,
		Nonce:     nonce,
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
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["email"] = user.Email
	data["product"] = product.Name
	err = s.EmailService.Send("payment", data)
	if err != nil {
		return response, err
	}

	return CreatePaymentResponse{
		PaymentID:   pay.ID,
		UserID:      user.ID,
		ProductID:   product.ID,
		Price:       amount.GetFloat(),
		Currency:    product.Currency,
		ProductName: product.Name,
		UserName:    user.Name,
		Token:       user.JWT,
		Success:     true,
	}, nil
}
