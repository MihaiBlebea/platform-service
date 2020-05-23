package payment

import (
	"time"
)

// Provider Payment interface
type Provider interface {
	Connect()
	Payment() (string, error)
}

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

// Details DTO
type Details struct {
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}

// Make a new transaction with the payment provider
func Make(provider Provider, repository Repository, details Details) (payment Payment, err error) {
	provider.Connect()
	_, err = provider.Payment()
	if err != nil {
		return payment, err
	}

	payment = Payment{
		UserID:      details.UserID,
		ProductID:   details.ProductID,
		PaymentCode: "abcd_payment",
		Price:       details.Price,
		Currency:    details.Currency,
		Created:     time.Now(),
	}

	paymentID, err := repository.Add(&payment)
	if err != nil {
		return payment, err
	}
	payment.ID = paymentID

	return payment, nil
}
