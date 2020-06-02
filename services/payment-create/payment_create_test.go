package paycreate

// func registerUser() (user *u.User, err error) {
// 	var (
// 		name     = "Serban"
// 		email    = "serban@gmail.com"
// 		password = "intrex"
// 		consent  = true
// 	)

// 	userRepository := *u.Repo(&c.MysqlConnection{
// 		Username: "admin",
// 		Password: "pass",
// 		Host:     "127.0.0.1",
// 		Port:     "3306",
// 		Database: "platform",
// 	})
// 	user = &u.User{
// 		Name:     name,
// 		Email:    email,
// 		Password: password,
// 		Consent:  consent,
// 	}
// 	id, err := userRepository.Add(user)
// 	if err != nil {
// 		return user, err
// 	}
// 	user.ID = id

// 	return user, nil
// }

// func addProduct() (product p.Product, err error) {
// 	var (
// 		name  = "Domo"
// 		code  = "d-mo"
// 		price = 49.99
// 	)

// 	productRepository := *p.Repo(&c.MysqlConnection{
// 		Username: "admin",
// 		Password: "pass",
// 		Host:     "127.0.0.1",
// 		Port:     "3306",
// 		Database: "platform",
// 	})

// 	product = p.Product{
// 		Code:     code,
// 		Name:     name,
// 		Price:    price,
// 		Currency: "GBP",
// 		Active:   true,
// 	}
// 	id, err := productRepository.Add(&product)
// 	if err != nil {
// 		return product, err
// 	}
// 	product.ID = id

// 	return product, nil
// }

// func createMysqlConnection() *c.MysqlConnection {
// 	return &c.MysqlConnection{
// 		Username: "admin",
// 		Password: "pass",
// 		Host:     "127.0.0.1",
// 		Port:     "3306",
// 		Database: "platform",
// 	}
// }

// func TestPaymentWithAuth(t *testing.T) {
// 	registerResponse, err := registerUser()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	product, err := addProduct()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	userRepository := *u.Repo(createMysqlConnection())
// 	productRepository := *p.Repo(createMysqlConnection())
// 	tokenRepository := *tkn.Repo(createMysqlConnection())
// 	paymentRepository := *payment.Repo(createMysqlConnection())

// 	paymentService := CreatePaymentService{
// 		userRepository,
// 		productRepository,
// 		tokenRepository,
// 		paymentRepository,
// 	}

// 	response, err := paymentService.ExecuteWithAuth(registerResponse.JWT, product.Code)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if response.ProductID != product.ID {
// 		t.Errorf("Expected product id %d, got %d", product.ID, response.ProductID)
// 	}
// 	if response.UserID != registerResponse.ID {
// 		t.Errorf("Expected user id %d, got %d", registerResponse.ID, response.UserID)
// 	}
// 	if response.Success != true {
// 		t.Errorf("Expected Success to be true")
// 	}

// 	t.Cleanup(func() {
// 		productRepository.Remove(&p.Product{ID: product.ID})

// 		user := u.User{ID: registerResponse.ID}
// 		userRepository.Remove(&user)

// 		paymentRepository.Remove(
// 			&payment.Payment{
// 				ID: response.PaymentID,
// 			},
// 		)

// 		tokenRepository.RemoveByUser(registerResponse.ID)
// 	})
// }

// func TestPaymentWithoutAuth(t *testing.T) {
// 	var (
// 		firstName = "John"
// 		lastName  = "Doe"
// 		email     = "john.doe@gmail.com"
// 	)
// 	product, err := addProduct()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	userRepository := *u.Repo(createMysqlConnection())
// 	productRepository := *p.Repo(createMysqlConnection())
// 	tokenRepository := *tkn.Repo(createMysqlConnection())
// 	paymentRepository := *payment.Repo(createMysqlConnection())

// 	paymentService := CreatePaymentService{
// 		userRepository,
// 		productRepository,
// 		tokenRepository,
// 		paymentRepository,
// 	}

// 	response, err := paymentService.Execute(
// 		product.Code,
// 		firstName,
// 		lastName,
// 		email,
// 	)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if response.ProductID != product.ID {
// 		t.Errorf("Expected product id %d, got %d", product.ID, response.ProductID)
// 	}
// 	if response.Success != true {
// 		t.Errorf("Expected Success to be true")
// 	}
// 	if response.UserName != fmt.Sprintf("%s %s", firstName, lastName) {
// 		t.Errorf("Expected Username %s, but got %s", fmt.Sprintf("%s %s", firstName, lastName), response.UserName)
// 	}

// 	t.Cleanup(func() {
// 		productRepository.Remove(&p.Product{ID: product.ID})

// 		user := u.User{ID: response.UserID}
// 		userRepository.Remove(&user)

// 		paymentRepository.Remove(
// 			&payment.Payment{
// 				ID: response.PaymentID,
// 			},
// 		)

// 		tokenRepository.RemoveByUser(response.UserID)
// 	})
// }
