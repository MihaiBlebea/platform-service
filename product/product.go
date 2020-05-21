package product

import "time"

// Product model
type Product struct {
	ID       int       `json:"id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Currency string    `json:"currency"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}
