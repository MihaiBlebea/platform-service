package paycreate

import (
	"fmt"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
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

	return &CreatePaymentService{
		userRepository,
		productRepository,
		tokenRepository,
		paymentRepository,
	}
}

// CreatePaymentService creates a payment
type CreatePaymentService struct {
	UserRepository    u.Repository
	ProductRepository p.Repository
	TokenRepository   t.Repository
	PaymentRepository payment.Repository
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
func (s *CreatePaymentService) ExecuteWithAuth(token, productCode string) (response CreatePaymentResponse, err error) {
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
	details := payment.Details{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     product.Price,
		Currency:  "GBP",
	}
	pay, err := payment.Make(
		&payment.BraintreeProvider{},
		s.PaymentRepository,
		details,
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

	return CreatePaymentResponse{
		PaymentID:   pay.ID,
		UserID:      user.ID,
		ProductID:   product.ID,
		Price:       product.Price,
		Currency:    product.Currency,
		ProductName: product.Name,
		UserName:    user.Name,
		Token:       user.JWT,
		Success:     true,
	}, nil
}

// Execute creates a payment when the user is not logged in the platform
func (s *CreatePaymentService) Execute(productCode, firstName, lastName, email string) (response CreatePaymentResponse, err error) {
	user, err := u.New(
		fmt.Sprintf("%s %s", firstName, lastName),
		email,
		u.GenerateRandomPassword(),
		true,
	)
	if err != nil {
		return response, err
	}

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
	details := payment.Details{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     product.Price,
		Currency:  "GBP",
	}
	pay, err := payment.Make(
		&payment.BraintreeProvider{},
		s.PaymentRepository,
		details,
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

	return CreatePaymentResponse{
		PaymentID:   pay.ID,
		UserID:      user.ID,
		ProductID:   product.ID,
		Price:       product.Price,
		Currency:    product.Currency,
		ProductName: product.Name,
		UserName:    user.Name,
		Token:       user.JWT,
		Success:     true,
	}, nil
}
