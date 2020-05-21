package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"github.com/MihaiBlebea/Wordpress/platform/product"
	"github.com/julienschmidt/httprouter"
)

func productGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	code := params.Get("code")
	if code == "" {
		log.Panic(errors.New("No product code supplied"))
	}

	productRepo := product.Repo(c.Mysql())
	product, count, err := productRepo.FindByCode(code)
	if err != nil {
		log.Panic(err)
	}
	if count == 0 {
		log.Panic(err)
	}

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
