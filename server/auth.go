package server

import (
	"encoding/json"
	"net/http"

	useremailconfirm "github.com/MihaiBlebea/purpletree/platform/services/user-email-confirm"
	usremailvalidate "github.com/MihaiBlebea/purpletree/platform/services/user-email-validate"
	usrlogin "github.com/MihaiBlebea/purpletree/platform/services/user-login"
	usrpassconfirm "github.com/MihaiBlebea/purpletree/platform/services/user-password-confirm"
	usrpassreset "github.com/MihaiBlebea/purpletree/platform/services/user-password-reset"
	useregister "github.com/MihaiBlebea/purpletree/platform/services/user-register"
	"github.com/julienschmidt/httprouter"
)

// Login to the platform
func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Email    string
		Password string
	}
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := usrlogin.New()
	response, err := service.Execute(body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Register to the platform, without confirming your account
func registerPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Name     string
		Email    string
		Password string
		Consent  bool
	}
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := useregister.New()
	response, err := service.Execute(body.Name, body.Email, body.Password, body.Consent)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Confirm your email and activate your account
func registerConfirmPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Code string
	}
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := useremailconfirm.New()
	response, err := service.Execute(body.Code)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Reset user's password
func passwordPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		ConfirmCode string
		Password    string
	}
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := usrpassreset.New()
	response, err := service.Execute(body.ConfirmCode, body.Password)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Validate user's email and redirect to change password page
func passwordGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	email := params.Get("email")
	if email == "" {
		http.Error(w, "No email supplied", 500)
		return
	}

	confirmEndpoint := params.Get("confirmEndpoint")
	if confirmEndpoint == "" {
		http.Error(w, "No confirmEndpoint supplied", 500)
		return
	}

	service := usrpassconfirm.New()
	response, err := service.Execute(confirmEndpoint, email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func emailValidateGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	email := params.Get("email")
	if email == "" {
		http.Error(w, "No emailsupplied", 500)
	}

	service := usremailvalidate.New()
	response, err := service.Execute(email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
