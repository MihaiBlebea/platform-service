package server

import (
	"encoding/json"
	"net/http"

	useremailconfirm "github.com/MihaiBlebea/Wordpress/platform/services/user-email-confirm"
	useregister "github.com/MihaiBlebea/Wordpress/platform/services/user-register"
	"github.com/MihaiBlebea/Wordpress/platform/services/usrlogin"
	"github.com/MihaiBlebea/Wordpress/platform/services/usrpassconfirm"
	"github.com/MihaiBlebea/Wordpress/platform/services/usrpassreset"
	"github.com/julienschmidt/httprouter"
)

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Email    string
		Password string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
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

func registerPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Name     string
		Email    string
		Password string
		Consent  bool
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
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

func registerConfirmPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Code string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
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

func passwordPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		ConfirmCode string
		Password    string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
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
