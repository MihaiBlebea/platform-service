package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/MihaiBlebea/Wordpress/platform/services/checkout"
	paycreate "github.com/MihaiBlebea/Wordpress/platform/services/payment-create"
)

func paymentPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		FirstName   string
		LastName    string
		Email       string
		Nonce       string
		PaymentType string
		Consent     bool
		ProductCode string
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

		response, err = service.ExecuteWithAuth(
			jwt,
			body.Nonce,
			body.PaymentType,
			body.ProductCode,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		response, err = service.Execute(
			body.FirstName,
			body.LastName,
			body.Email,
			body.Nonce,
			body.PaymentType,
			body.ProductCode,
		)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

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
