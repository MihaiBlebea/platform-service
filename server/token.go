package server

import (
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/Wordpress/platform/services/tknactiv"
	"github.com/MihaiBlebea/Wordpress/platform/services/tkndeactiv"
	"github.com/MihaiBlebea/Wordpress/platform/services/tknvalid"
	"github.com/MihaiBlebea/Wordpress/platform/services/usrtkns"
	"github.com/julienschmidt/httprouter"
)

func activateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Token string
		Host  string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := tknactiv.New()
	response, err := service.Execute(body.Token, body.Host)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func deactivateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Token string
		Host  string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := tkndeactiv.New()
	response, err := service.Execute(body.Token, body.Host)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func vaidateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Token string
		Host  string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := tknvalid.New()
	response := service.Execute(body.Token, body.Host)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func tokensGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if authenticate(r) == false {
		http.Error(w, "Not authenticated", 401)
		return
	}

	user, err := authenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := usrtkns.New()
	response, err := service.Execute(user.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
