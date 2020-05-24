package prodget

import (
	"testing"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	p "github.com/MihaiBlebea/Wordpress/platform/product"
)

func TestProductCannotBeFetched(t *testing.T) {
	var (
		code = "d-mo"
	)

	productRepository := *p.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})
	service := GetProductService{
		productRepository,
	}
	_, err := service.Execute(code)
	if err == nil {
		t.Error("Should throw an error")
	}
}

func TestProductCanBeFetched(t *testing.T) {
	var (
		name  = "Domo"
		code  = "d-mo"
		price = 49.99
	)

	productRepository := *p.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	productRepository.Add(
		&p.Product{
			Code:     code,
			Name:     name,
			Price:    price,
			Currency: "GBP",
			Active:   true,
		},
	)
	service := GetProductService{
		productRepository,
	}
	response, err := service.Execute(code)
	if err != nil {
		t.Error(err)
	}
	if response.Name != name {
		t.Errorf("Expected name %s, got %s", name, response.Name)
	}
	if response.Code != code {
		t.Errorf("Expected code %s, got %s", code, response.Code)
	}
	if response.Price != price {
		t.Errorf("Expected price %d, got %d", int(price), int(response.Price))
	}

	t.Cleanup(func() {
		productRepository.Remove(&p.Product{ID: response.ID})
	})
}
