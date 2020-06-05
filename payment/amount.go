package payment

// Amount of money
type Amount struct {
	amount   int
	currency string
}

// NewAmount creates a new Amount object
func NewAmount(amount int, currency string) *Amount {
	return &Amount{amount, currency}
}

// GetFloat returns the amount as a float with 2 decimals
func (a *Amount) GetFloat() float64 {
	return float64(a.amount / 100)
}

// GetInt returns the amount as an int
func (a *Amount) GetInt() int {
	return a.amount
}

// GetCurrency returns the currency of the amount
func (a *Amount) GetCurrency() string {
	return a.currency
}

// FloatWithTVA returns amount as float with added TVA of 20%
func (a *Amount) FloatWithTVA() float64 {
	amount := a.GetFloat()
	return amount + amount*0.2
}

// IntWithTVA returns amount as int with added TVA of 20%
func (a *Amount) IntWithTVA() int {
	total := a.FloatWithTVA()
	return int(total * 100)
}
