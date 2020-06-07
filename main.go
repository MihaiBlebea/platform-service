package main

import (
	"log"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"github.com/MihaiBlebea/Wordpress/platform/discount"
	"github.com/MihaiBlebea/Wordpress/platform/payment"
	"github.com/MihaiBlebea/Wordpress/platform/product"
	"github.com/MihaiBlebea/Wordpress/platform/server"
	"github.com/MihaiBlebea/Wordpress/platform/user"
	"github.com/MihaiBlebea/Wordpress/platform/user/token"
)

const serverPort = ":8001"

func init() {
	userRepo := user.Repo(c.Mysql())
	tokenRepo := token.Repo(c.Mysql())
	productRepo := product.Repo(c.Mysql())
	paymentRepo := payment.Repo(c.Mysql())
	discountRepo := discount.Repo(c.Mysql())

	err := userRepo.Migrate()
	if err != nil {
		log.Panic(err)
	}
	err = tokenRepo.Migrate()
	if err != nil {
		log.Panic(err)
	}
	err = productRepo.Migrate()
	if err != nil {
		log.Panic(err)
	}
	err = paymentRepo.Migrate()
	if err != nil {
		log.Panic(err)
	}
	err = discountRepo.Migrate()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	server.Serve(serverPort)
}
