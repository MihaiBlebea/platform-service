package payment

import (
	"errors"
	"time"
)

// Service is a payment service
type Service struct {
}

// Params DTO
type Params struct {
	UserID    int
	ProductID int
	Price     Price
	Nonce     string
}

// Pay a new transaction with the payment provider
func (s *Service) Pay(repository Repository, params Params) (payment Payment, err error) {
	provider := BraintreeProvider{}

	// Validate the params
	err = s.validate(params)
	if err != nil {
		return payment, err
	}

	// make a new payment with nonce
	transaction, err := provider.paymentWithNonce(params.Nonce, params.Price.WithTVA())
	if err != nil {
		return payment, err
	}

	// Save the payment information in the db
	payment = Payment{
		UserID:      params.UserID,
		ProductID:   params.ProductID,
		PaymentCode: transaction.Id,
		Price:       params.Price.WithTVA(),
		Currency:    "GBP",
		Created:     time.Now(),
	}

	paymentID, err := repository.Add(&payment)
	if err != nil {
		return payment, err
	}
	payment.ID = paymentID

	// Return response
	return payment, nil
}

// GetClientToken returns a token that can be used in frontend for getting a nonce
func (s *Service) GetClientToken() (token string, err error) {
	provider := BraintreeProvider{}

	return provider.getClientToken()
}

func (s *Service) validate(params Params) error {
	if params.UserID == 0 {
		return errors.New("UserID is not set")
	}

	if params.Price.GetAmount() == 0 {
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
