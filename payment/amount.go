package payment

// Price of money
type Price struct {
	amount   float64
	currency string
}

// NewPrice creates a new Price object
func NewPrice(amount float64, currency string) *Price {
	return &Price{amount, currency}
}

// ApplyDiscount applies a discount to a price
func (p *Price) ApplyDiscount(percentage float64) {
	p.amount = p.amount + p.amount*percentage
}

// GetAmount returns the amount as a float with 2 decimals
func (p *Price) GetAmount() float64 {
	return p.amount
}

// GetCurrency returns the currency of the amount
func (p *Price) GetCurrency() string {
	return p.currency
}

// WithTVA returns amount as float with added TVA of 20%
func (p *Price) WithTVA() float64 {
	return p.amount + p.amount*0.2
}
