package payment

import (
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
