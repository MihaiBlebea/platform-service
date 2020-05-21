package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MihaiBlebea/Wordpress/platform/user/token"
	"github.com/julienschmidt/httprouter"
)

func deactivateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Body struct {
		Token string
		Host  string
	}
	decoder := json.NewDecoder(r.Body)
	var body Body
	err := decoder.Decode(&body)
	if err != nil {
		log.Panic(err)
	}
	response, err := token.Deactivate(body.Token, body.Host)
	if err != nil {
		log.Panic(err)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
