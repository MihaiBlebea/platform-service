package server

import (
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/Wordpress/platform/services/usregister"
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

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	service := usregister.New()
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

func passwordPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	service := usrpassreset.New()
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

func passwordGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	email := params.Get("email")
	if email == "" {
		http.Error(w, "No email supplied", 500)
		return
	}

	service := usrpassconfirm.New()
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
