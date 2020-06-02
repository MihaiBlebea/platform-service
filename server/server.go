package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	u "github.com/MihaiBlebea/Wordpress/platform/user"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Serve initializes the http server
func Serve(port string) {
	router := httprouter.New()

	router.GET("/", indexHandler)

	router.POST("/user/register", registerPostHandler)
	router.POST("/user/register-confirm-email", registerConfirmPostHandler)
	router.POST("/user/login", loginHandler)

	router.GET("/user/password/reset", passwordGetHandler)
	router.POST("/user/password/reset", passwordPostHandler)

	router.GET("/user/tokens", tokensGetHandler)
	router.POST("/user/token/validate", vaidateTokenHandler)
	router.POST("/user/token/activate", activateTokenHandler)
	router.POST("/user/token/deactivate", deactivateTokenHandler)

	router.GET("/product", productGetHandler)
	router.POST("/payment", paymentPostHandler)
	router.GET("/payment/checkout", paymentCheckoutGetHandler)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := corsMiddleware.Handler(router)

	// handler := cors.Default().Handler(router)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Panic(err)
	}
}

func authenticate(r *http.Request) bool {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return false
	}
	token := strings.TrimSpace(strings.Split(authorization, "Bearer")[1])
	found, _, err := u.Authenticate(token)
	if err != nil {
		return false
	}
	return found
}

func authenticatedUser(r *http.Request) (user *u.User, err error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return user, errors.New("Could not find JWT token in request")
	}
	token := strings.TrimSpace(strings.Split(authorization, "Bearer")[1])
	_, user, err = u.Authenticate(token)
	if err != nil {
		return user, err
	}
	return user, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	health := make(map[string]string)
	health["status"] = "OK"

	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
