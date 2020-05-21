package server

import (
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/Wordpress/platform/user"
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
	}

	response, err := user.Login(body.Email, body.Password)
	if err != nil || response.Success == false {
		w.WriteHeader(403)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
