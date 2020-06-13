package server

import (
	"encoding/json"
	"net/http"

	prodget "github.com/MihaiBlebea/purpletree/platform/services/product-get"
	"github.com/julienschmidt/httprouter"
)

func productGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	code := params.Get("code")
	if code == "" {
		http.Error(w, "No product code supplied", 500)
	}

	service := prodget.New()
	response, err := service.Execute(code)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
