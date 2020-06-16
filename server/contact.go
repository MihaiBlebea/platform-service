package server

import (
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/purpletree/platform/services/contact"
	"github.com/julienschmidt/httprouter"
)

func contactPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Name    string
		Email   string
		Subject string
		Message string
	}
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	service := contact.New()
	response, err := service.Execute(body.Name, body.Email, body.Subject, body.Message)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
