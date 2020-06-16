package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/purpletree/platform/server/limiter"
	usrauth "github.com/MihaiBlebea/purpletree/platform/services/user-authenticate"
	u "github.com/MihaiBlebea/purpletree/platform/user"

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
	router.GET("/user/email-validate", emailValidateGetHandler)

	router.GET("/user/tokens", tokensGetHandler)
	router.POST("/user/token/validate", vaidateTokenHandler)
	router.POST("/user/token/activate", activateTokenHandler)
	router.POST("/user/token/deactivate", deactivateTokenHandler)

	router.GET("/product", productGetHandler)

	router.POST("/payment", paymentPostHandler)
	router.GET("/payment/checkout", paymentCheckoutGetHandler)

	router.POST("/contact", contactPostHandler)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := corsMiddleware.Handler(router)

	handler = limiter.Limit(handler)

	handler = auth(handler)

	// handler := cors.Default().Handler(router)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Panic(err)
	}
}

// auth is a middleware that checks if the current route
// is part of the protected list of routes,
// if so, then it requires the request to have a valid JWT or it's returned with error
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protected := []string{"/user/tokens"}

		contains := func(s []string, e string) bool {
			for _, a := range s {
				if a == e {
					return true
				}
			}
			return false
		}

		if contains(protected, r.URL.Path) == false {
			next.ServeHTTP(w, r)
			return
		}

		_, err := getUserFromJWT(r)
		if err != nil {
			http.Error(w, "No Authorization", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUserFromJWT(r *http.Request) (user *u.User, err error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return user, errors.New("No auth token provided")
	}

	jwt := strings.TrimSpace(strings.Split(authorization, "Bearer")[1])

	service := usrauth.New()
	isValid, user, err := service.Execute(jwt)
	if err != nil {
		return user, err
	}
	if isValid == false {
		return user, errors.New("No auth token provided")
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
