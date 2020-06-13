package server

import (
	"encoding/json"
	"net/http"

	tknactiv "github.com/MihaiBlebea/purpletree/platform/services/token-activate"
	tkndeactiv "github.com/MihaiBlebea/purpletree/platform/services/token-deactivate"
	tknvalid "github.com/MihaiBlebea/purpletree/platform/services/token-validate"
	usrtkns "github.com/MihaiBlebea/purpletree/platform/services/user-tokens"
	"github.com/julienschmidt/httprouter"
)

func activateTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// type Body struct {
	// 	Token string
	// 	Host  string
	// }
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Token string
		Host  string
	}
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
	// type Body struct {
	// 	Token string
	// 	Host  string
	// }
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Token string
		Host  string
	}
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
	// type Body struct {
	// 	Token string
	// 	Host  string
	// }
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Token string
		Host  string
	}
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
