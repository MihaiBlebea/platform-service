package checkout

import (
	"testing"
	"time"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	d "github.com/MihaiBlebea/purpletree/platform/discount"
	pay "github.com/MihaiBlebea/purpletree/platform/payment"
	p "github.com/MihaiBlebea/purpletree/platform/product"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

func newFakeService() Service {
	productRepository := *p.Repo(c.Mysql())
	discountRepository := *d.Repo(c.Mysql())
	userRepository := *u.Repo(c.Mysql())
	paymentService := &pay.Service{}

	return Service{
		productRepository,
		discountRepository,
		userRepository,
		paymentService,
	}
}

func newFakeUser() (u.User, error) {
	userRepository := *u.Repo(c.Mysql())
	user := u.User{
		Name:     "Mihai",
		Email:    "mihai@gmail.com",
		Password: "abcd1234",
		Active:   true,
		Consent:  true,
	}

	user.GenerateJWT()
	userID, err := userRepository.Add(&user)
	if err != nil {
		return user, err
	}
	user.ID = userID

	return user, nil
}

func removeFakeUser(user u.User) error {
	userRepository := *u.Repo(c.Mysql())
	err := userRepository.Remove(&user)

	return err
}

func newFakeProduct() (p.Product, error) {
	productRepository := *p.Repo(c.Mysql())

	product := p.Product{
		Code:     "abcd",
		Name:     "dust buster",
		Price:    79.99,
		Currency: "GBP",
		Active:   true,
	}

	productID, err := productRepository.Add(&product)
	if err != nil {
		return product, err
	}

	product.ID = productID

	return product, nil
}

func removeFakeProduct(product p.Product) error {
	productRepository := *p.Repo(c.Mysql())
	err := productRepository.Remove(&product)

	return err
}

func newFakeDiscount(productID int) (d.Discount, error) {
	discountRepository := *d.Repo(c.Mysql())

	discount := d.Discount{
		ProductID:  productID,
		Code:       "DISCOUNT_CODE_1",
		Percentage: 0.5,
		Active:     true,
		Expires:    time.Now(),
	}

	discountID, err := discountRepository.Add(&discount)
	if err != nil {
		return discount, err
	}

	discount.ID = discountID

	return discount, nil
}

func removeFakeDiscount(discount d.Discount) error {
	discountRepository := *d.Repo(c.Mysql())
	err := discountRepository.Remove(&discount)

	return err
}

func TestCheckoutWithUser(t *testing.T) {
	service := newFakeService()

	user, err := newFakeUser()
	if err != nil {
		t.Fatal(err)
	}

	product, err := newFakeProduct()
	if err != nil {
		t.Fatal(err)
	}

	discount, err := newFakeDiscount(product.ID)
	if err != nil {
		t.Fatal(err)
	}

	response, err := service.Execute(Request{
		ProductCode:  product.Code,
		JWT:          user.JWT,
		DiscountCode: discount.Code,
	})
	if err != nil {
		t.Fatal(err)
	}

	if response.UserEmail != user.Email {
		t.Errorf("Expected response user email %s, got %s", user.Email, response.UserEmail)
	}

	t.Cleanup(func() {
		removeFakeUser(user)
		removeFakeProduct(product)
		removeFakeDiscount(discount)
	})
}
