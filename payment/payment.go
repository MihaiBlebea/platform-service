package payment

import (
	"errors"
	"time"
)

// Payment model
type Payment struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProductID   int       `json:"product_id"`
	PaymentCode string    `json:"payment_code"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Params DTO
type Params struct {
	UserID    int
	ProductID int
	Price     float64
	Nonce     string
}

// Pay a new transaction with the payment provider
func Pay(repository Repository, params Params) (payment Payment, err error) {
	provider := BraintreeProvider{}

	err = validate(params)
	if err != nil {
		return payment, err
	}

	_, err = provider.paymentWithNonce(params.Nonce, int64(params.Price*100))
	if err != nil {
		return payment, err
	}

	payment = Payment{
		UserID:      params.UserID,
		ProductID:   params.ProductID,
		PaymentCode: "abcd_payment",
		Price:       params.Price,
		Currency:    "GBP",
		Created:     time.Now(),
	}

	paymentID, err := repository.Add(&payment)
	if err != nil {
		return payment, err
	}
	payment.ID = paymentID

	return payment, nil
}

// GetClientToken returns a token that can be used in frontend for getting a nonce
func GetClientToken() (token string, err error) {
	provider := BraintreeProvider{}

	return provider.getClientToken()
}

func validate(params Params) error {
	if params.UserID == 0 {
		return errors.New("UserID is not set")
	}

	if params.Price == 0 {
		return errors.New("Price is not set")
	}

	if params.ProductID == 0 {
		return errors.New("ProductID is not set")
	}

	if params.Nonce == "" {
		return errors.New("Nonce is not set")
	}
	return nil
}
