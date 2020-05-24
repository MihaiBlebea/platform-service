package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/MihaiBlebea/Wordpress/platform/services/paycreate"
)

type paymentRequestBody struct {
	FirstName   string
	LastName    string
	Email       string
	CardName    string
	CardNumber  string
	ExpireMonth string
	ExpireYear  string
	CVV         string
	Consent     bool
	ProductCode string
}

func paymentPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body paymentRequestBody
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

		response, err = service.ExecuteWithAuth(jwt, body.ProductCode)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		response, err = service.Execute(body.ProductCode, body.FirstName, body.LastName, body.Email)
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
