package prodget

import (
	"fmt"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	p "github.com/MihaiBlebea/Wordpress/platform/product"
)

// New returns a new GetProductService
func New() *GetProductService {
	productRepository := *p.Repo(c.Mysql())
	return &GetProductService{productRepository}
}

// GetProductService returns a product from the database
type GetProductService struct {
	productRepository p.Repository
}

// GetProductResponse is the response struct for GetProductService
type GetProductResponse struct {
	ID       int     `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// Execute runs the RegisterUserService
func (s *GetProductService) Execute(code string) (response GetProductResponse, err error) {
	product, count, err := s.productRepository.FindByCode(code)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("Could not find product with code %s", code)
	}

	return GetProductResponse{
		ID:       product.ID,
		Code:     product.Code,
		Name:     product.Name,
		Price:    product.Price,
		Currency: product.Currency,
	}, nil
}
