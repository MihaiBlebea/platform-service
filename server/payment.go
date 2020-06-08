package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/MihaiBlebea/Wordpress/platform/services/checkout"
	paycreate "github.com/MihaiBlebea/Wordpress/platform/services/payment-create"
	paydiscountvalid "github.com/MihaiBlebea/Wordpress/platform/services/payment-discount-validate"
)

// Post a payment request with or without auth. If no auth code, then create account for user
func paymentPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		FirstName    string
		LastName     string
		Email        string
		Password     string
		Nonce        string
		PaymentType  string
		Consent      bool
		ProductCode  string
		DiscountCode string
	}

	decoder := json.NewDecoder(r.Body)
	var body Body
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := paycreate.New()
	var response paycreate.CreatePaymentResponse

	// Fetch user account from jwt token or create a new user if no token is supplied
	jwt := r.Header.Get("Authorization")
	if jwt != "" {
		jwt = strings.TrimSpace(strings.Split(jwt, "Bearer")[1])
	}

	request := paycreate.CreatePaymentRequest{
		Token:        jwt,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		Password:     body.Password,
		Nonce:        body.Nonce,
		PaymentType:  body.PaymentType,
		ProductCode:  body.ProductCode,
		DiscountCode: body.DiscountCode,
	}
	response, err = service.Execute(request)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Get all the information to render the checkout page
func paymentCheckoutGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	code := params.Get("code")
	if code == "" {
		http.Error(w, "No product code supplied", 500)
	}

	jwt := r.Header.Get("Authorization")
	if jwt != "" {
		jwt = strings.TrimSpace(strings.Split(jwt, "Bearer")[1])
	}

	service := checkout.New()
	response, err := service.Execute(code, jwt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func discountValidateGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	code := params.Get("code")
	if code == "" {
		http.Error(w, "No discount code supplied", 500)
	}
	product := params.Get("product")
	if code == "" {
		http.Error(w, "No product code supplied", 500)
	}
	productID, err := strconv.Atoi(product)
	if err != nil {
		http.Error(w, "Invalid product id", 500)
	}

	service := paydiscountvalid.New()
	response, err := service.Execute(code, productID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
