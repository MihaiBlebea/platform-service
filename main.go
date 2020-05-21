package main

import (
	"log"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"github.com/MihaiBlebea/Wordpress/platform/payment"
	"github.com/MihaiBlebea/Wordpress/platform/product"
	"github.com/MihaiBlebea/Wordpress/platform/server"
	"github.com/MihaiBlebea/Wordpress/platform/user"
	"github.com/MihaiBlebea/Wordpress/platform/user/token"
)

func main() {
	userRepo := user.Repo(c.Mysql())
	tokenRepo := token.Repo(c.Mysql())
	productRepo := product.Repo(c.Mysql())
	paymentRepo := payment.Repo(c.Mysql())

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

	server.Serve()
}
