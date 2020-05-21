package payment

import (
	"context"
	"fmt"

	"github.com/braintree-go/braintree-go"
)

const merchantID = "6ptmhsvtnc8jdtc4"
const publicKey = "x985fgd2ykxb4s5s"
const privateKey = "b85a60351560cdcaeabafe713c69ec92"

// https://sandbox.braintreegateway.com/merchants/6ptmhsvtnc8jdtc4/transactions/advanced_search

// BraintreeProvider payment provider
type BraintreeProvider struct {
	Client braintree.Braintree
}

// Connect returns a braintree client object
func (p *BraintreeProvider) Connect() {
	p.Client = *braintree.New(
		braintree.Sandbox,
		merchantID,
		publicKey,
		privateKey,
	)
}

// Payment makes a new transaction
func (p *BraintreeProvider) Payment() (string, error) {
	ctx := context.Background()
	transaction, err := p.Client.Transaction().Create(ctx, &braintree.TransactionRequest{
		Type:   "sale",
		Amount: braintree.NewDecimal(100, 2), // 100 cents
		CreditCard: &braintree.CreditCard{
			Number:         "4111111111111111",
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Println(transaction)
	return "", nil
}
