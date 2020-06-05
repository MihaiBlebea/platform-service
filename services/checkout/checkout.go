package checkout

import (
	"errors"
	"fmt"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"github.com/MihaiBlebea/Wordpress/platform/payment"
	pay "github.com/MihaiBlebea/Wordpress/platform/payment"
	p "github.com/MihaiBlebea/Wordpress/platform/product"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
)

// New returns a new GetProductService
func New() *Service {
	productRepository := *p.Repo(c.Mysql())
	userRepository := *u.Repo(c.Mysql())
	paymentService := &pay.Service{}

	return &Service{
		productRepository,
		userRepository,
		paymentService,
	}
}

// Service is a service
type Service struct {
	ProductRepository p.Repository
	UserRepository    u.Repository
	PaymentService    *pay.Service
}

// Response is a response for CheckoutService
type Response struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Price       float64     `json:"price"`
	PriceTotal  float64     `json:"price_total"`
	UserEmail   interface{} `json:"email"`
	ClientToken string      `json:"client_token"`
}

// Execute runs the Service
func (s *Service) Execute(code, jwtToken string) (response Response, err error) {
	// Find product with code
	product, count, err := s.ProductRepository.FindByCode(code)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("Could not find product with code %s", code)
	}

	// Get the Braintree client token
	token, err := s.PaymentService.GetClientToken()
	if err != nil {
		return response, err
	}

	// Get the user details
	if jwtToken != "" {
		user, count, err := s.UserRepository.FindByJWT(jwtToken)
		if err != nil {
			return response, err
		}
		if count == 0 {
			return response, errors.New("Could not find user")
		}

		response.UserEmail = user.Email
	} else {
		response.UserEmail = nil
	}

	amount := payment.NewAmount(product.Price, product.Currency)

	response.Name = product.Name
	response.Price = amount.GetFloat()
	response.PriceTotal = amount.FloatWithTVA()
	response.Code = product.Code
	response.ClientToken = token

	return response, nil
}
