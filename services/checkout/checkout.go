package checkout

import (
	"errors"
	"fmt"
	"strings"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	d "github.com/MihaiBlebea/purpletree/platform/discount"
	"github.com/MihaiBlebea/purpletree/platform/payment"
	pay "github.com/MihaiBlebea/purpletree/platform/payment"
	p "github.com/MihaiBlebea/purpletree/platform/product"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New returns a new GetProductService
func New() *Service {
	productRepository := *p.Repo(c.Mysql())
	discountRepository := *d.Repo(c.Mysql())
	userRepository := *u.Repo(c.Mysql())
	paymentService := &pay.Service{}

	return &Service{
		productRepository,
		discountRepository,
		userRepository,
		paymentService,
	}
}

// Service is a service
type Service struct {
	ProductRepository  p.Repository
	DiscountRepository d.Repository
	UserRepository     u.Repository
	PaymentService     *pay.Service
}

// Response is a response for CheckoutService
type Response struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Price       Price       `json:"price"`
	Discount    Discount    `json:"discount"`
	UserEmail   interface{} `json:"email"`
	ClientToken string      `json:"client_token"`
}

// Price type as part of response
type Price struct {
	Original   float64 `json:"original"`
	Discounted float64 `json:"discounted"`
	WithTVA    float64 `json:"total"`
}

// Discount type as part of response
type Discount struct {
	Percentage float64 `json:"percentage"`
	Code       string  `json:"code"`
}

// Request is a request for CheckoutService
type Request struct {
	ProductCode  string
	JWT          string
	DiscountCode string
}

// Execute runs the Service
func (s *Service) Execute(request Request) (response Response, err error) {
	// Find product with code
	product, count, err := s.ProductRepository.FindByCode(request.ProductCode)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("Could not find product with code %s", request.ProductCode)
	}

	price := payment.NewPrice(product.Price, product.Currency)

	// Set the default price and TVA, will be overwriten if discount code is valid
	response.Price = Price{
		Original: price.GetAmount(),
		WithTVA:  price.WithTVA(),
	}

	// Find discount by code and apply it if found
	if request.DiscountCode != "" {
		discount, count, err := s.DiscountRepository.FindByCode(strings.ToUpper(request.DiscountCode))
		if err == nil && count > 0 && discount.IsValid() == true && discount.ProductID == product.ID {
			price.ApplyDiscount(discount.Percentage)

			response.Price = Price{
				Original:   price.GetOriginalAmount(),
				Discounted: price.GetDiscountedAmount(),
				WithTVA:    price.WithTVA(),
			}

			response.Discount = Discount{
				Percentage: discount.Percentage,
				Code:       strings.ToUpper(discount.Code),
			}
		}
	}

	// Get the Braintree client token
	token, err := s.PaymentService.GetClientToken()
	if err != nil {
		return response, err
	}

	// Get the user details
	if request.JWT != "" {
		user, count, err := s.UserRepository.FindByJWT(request.JWT)
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

	response.Name = product.Name
	response.Code = product.Code
	response.ClientToken = token

	return response, nil
}
