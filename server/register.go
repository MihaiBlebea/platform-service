package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MihaiBlebea/Wordpress/platform/user"
	"github.com/julienschmidt/httprouter"
)

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
		log.Panic(err)
	}

	response, err := user.Register(body.Name, body.Email, body.Password, body.Consent)
	if err != nil {
		log.Panic(err)
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
