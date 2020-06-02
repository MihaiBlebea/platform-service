package payment

import (
	"context"
	"fmt"
	"os"

	"github.com/braintree-go/braintree-go"
)

// https://sandbox.braintreegateway.com/merchants/6ptmhsvtnc8jdtc4/transactions/advanced_search

// BraintreeProvider payment provider
type BraintreeProvider struct {
	Client braintree.Braintree
}

// Connect returns a braintree client object
func (p *BraintreeProvider) Connect() *braintree.Braintree {
	return braintree.New(
		braintree.Sandbox,
		os.Getenv("BRAINTREE_MERCHANT_ID"),
		os.Getenv("BRAINTREE_PUBLIC_KEY"),
		os.Getenv("BRAINTREE_PRIVATE_KEY"),
	)
}

// func (p *BraintreeProvider) payment() (string, error) {
// 	ctx := context.Background()
// 	transaction, err := p.Client.Transaction().Create(ctx, &braintree.TransactionRequest{
// 		Type:   "sale",
// 		Amount: braintree.NewDecimal(100, 2), // 100 cents
// 		CreditCard: &braintree.CreditCard{
// 			Number:         "4111111111111111",
// 			ExpirationDate: "05/14",
// 		},
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println(transaction)
// 	return "", nil
// }

func (p *BraintreeProvider) paymentWithNonce(nonce string, amount int64) (result string, err error) {
	ctx := context.Background()
	client := p.Connect()

	_, err = client.PaymentMethodNonce().Find(ctx, nonce)
	if err != nil {
		return result, err
	}

	tx := &braintree.TransactionRequest{
		Type:               "sale",
		Amount:             braintree.NewDecimal(amount, 2),
		PaymentMethodNonce: nonce,
		Options: &braintree.TransactionOptions{
			ThreeDSecure: &braintree.TransactionOptionsThreeDSecureRequest{Required: false},
		},
	}

	txn, err := client.Transaction().Create(ctx, tx)
	if err != nil {
		return result, err
	}

	fmt.Println(txn)

	return result, nil
}

func (p *BraintreeProvider) getClientToken() (token string, err error) {
	ctx := context.Background()
	client := p.Connect()

	token, err = client.ClientToken().Generate(ctx)
	if err != nil {
		return token, err
	}

	return token, nil
}
