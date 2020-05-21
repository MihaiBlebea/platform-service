package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"github.com/MihaiBlebea/Wordpress/platform/email"
	"github.com/MihaiBlebea/Wordpress/platform/payment"
	prod "github.com/MihaiBlebea/Wordpress/platform/product"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
	"github.com/MihaiBlebea/Wordpress/platform/user/token"
	"github.com/julienschmidt/httprouter"
)

type paymentRequestBody struct {
	FirstName   string
	LastName    string
	Email       string
	CardName    string
	CardNumber  string
	ExpireMonth string
	ExpireYear  string
	CVV         string
	Consent     bool
	ProductCode string
}

func paymentPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var body paymentRequestBody
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Fetch user account from jwt token or create a new user if no token is supplied
	authorization := r.Header.Get("Authorization")
	fmt.Println("JWT", authorization)

	user, err := buildUser(authorization, body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("USER", user.ID)

	// Get product information from db
	productRepo := prod.Repo(c.Mysql())
	product, _, err := productRepo.FindByCode(body.ProductCode)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Make payment with payment provider
	details := payment.Details{
		UserID:    user.ID,
		ProductID: product.ID,
		Price:     product.Price,
		Currency:  "GBP",
	}
	pay, err := payment.Make(&payment.BraintreeProvider{}, details)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Generate activation tokens in user account
	tokens, err := token.NewCount(user.ID, user.Email, 3)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tokenRepo := token.Repo(c.Mysql())
	for _, token := range *tokens {
		_, err = tokenRepo.Add(&token)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// Send email to user with account information and payment confirmation
	type Data struct {
		Name string
	}
	emailDetails := email.Details{
		TemplateName:  "payment",
		Data:          Data{user.Name},
		Subject:       "Where is my money",
		ReceiverName:  user.Name,
		ReceiverEmail: user.Email,
	}
	_, err = email.Send(&email.SendInBlueProvider{}, emailDetails)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create response and return it
	type Response struct {
		Payment *payment.Payment `json:"payment"`
		User    *u.User          `json:"user"`
	}
	response := Response{pay, user}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func buildUser(header string, body paymentRequestBody) (*u.User, error) {
	userRepo := u.Repo(c.Mysql())

	if header != "" {
		jwt := strings.TrimSpace(strings.Split(header, "Bearer")[1])
		user, count, _ := userRepo.FindByJWT(jwt)
		if count > 0 {
			return user, nil
		}
	}

	user, err := u.New(
		fmt.Sprintf("%s %s", body.FirstName, body.LastName),
		body.Email,
		generateRandomPassword(),
		true,
	)
	if err != nil {
		return &u.User{}, err
	}
	userID, err := userRepo.Add(user)
	if err != nil {
		return &u.User{}, err
	}
	user.ID = userID

	return user, nil
}

func generateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
