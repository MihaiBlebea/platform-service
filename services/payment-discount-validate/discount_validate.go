package paydiscountvalid

import (
	"errors"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	d "github.com/MihaiBlebea/Wordpress/platform/discount"
)

// New returns a discount service
func New() *ValidateDiscountService {
	discountRepository := d.Repo(c.Mysql())
	return &ValidateDiscountService{discountRepository}
}

// ValidateDiscountService validates a discount code
type ValidateDiscountService struct {
	DiscountRepository *d.Repository
}

// ValidateDiscountResponse response for ValidateDiscountService
type ValidateDiscountResponse struct {
	Success    bool    `json:"success"`
	Percentage float64 `json:"percentage"`
}

// Execute runs the ValidateDiscountService service
func (s *ValidateDiscountService) Execute(code string, productID int) (response ValidateDiscountResponse, err error) {
	discount, count, err := s.DiscountRepository.FindByCode(code)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, err
	}
	if discount.IsValid() == false {
		return response, errors.New("Discount code is either expired or inactive")
	}
	if discount.ProductID != productID {
		return response, errors.New("Discount code is not valid for this product")
	}

	return ValidateDiscountResponse{true, discount.Percentage}, nil
}
