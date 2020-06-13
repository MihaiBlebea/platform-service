package payment

import "math"

const tva = 0.2

// Price of money
type Price struct {
	amount           float64
	originalAmount   float64
	discountedAmount float64
	currency         string
}

// NewPrice creates a new Price object
func NewPrice(amount float64, currency string) *Price {
	return &Price{
		amount:   amount,
		currency: currency,
	}
}

// ApplyDiscount applies a discount to a price
func (p *Price) ApplyDiscount(percentage float64) {
	p.originalAmount = p.amount
	p.discountedAmount = p.amount * percentage
	p.amount = p.amount - p.discountedAmount
}

// GetAmount returns the amount as a float with 2 decimals
func (p *Price) GetAmount() float64 {
	return toTwoDecimals(p.amount)
}

// GetOriginalAmount returns the original amount before applying discount as flaot with 2 decimals
func (p *Price) GetOriginalAmount() float64 {
	return toTwoDecimals(p.originalAmount)
}

// GetDiscountedAmount returns the discounted amount as flaot with 2 decimals
func (p *Price) GetDiscountedAmount() float64 {
	return toTwoDecimals(p.discountedAmount)
}

// GetCurrency returns the currency of the amount
func (p *Price) GetCurrency() string {
	return p.currency
}

// WithTVA returns amount as float with added TVA of 20%
func (p *Price) WithTVA() float64 {
	return toTwoDecimals(p.amount + p.amount*tva)
}

func toTwoDecimals(value float64) float64 {
	return math.Ceil(value*100) / 100
}
