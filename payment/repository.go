package payment

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"

	// Mysql driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// Repository struct
type Repository struct {
	Connection c.Connection
}

// Migrate generates a users table
func (r *Repository) Migrate() error {
	client, err := r.Connection.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	tableSchema := `CREATE TABLE IF NOT EXISTS payments(
		id INT NOT NULL AUTO_INCREMENT,
		user_id INT(100),
		product_id INT(100),
		payment_code VARCHAR(100),
		price FLOAT(6),
		currency VARCHAR(4),
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id));`
	stmt, err := client.Prepare(tableSchema)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// Add adds a User model to the db
func (r *Repository) Add(payment *Payment) (int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare("INSERT INTO payments (user_id, product_id, payment_code, price, currency, created) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(payment.UserID, payment.ProductID, payment.PaymentCode, payment.Price, payment.Currency, payment.Created)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// FindByProductID returns a user by id
func (r *Repository) FindByProductID(productID int) (*[]Payment, int, error) {
	payments, count, err := r.findBy("SELECT * FROM payments WHERE product_id = ?", productID)
	if err != nil {
		return &[]Payment{}, 0, err
	}
	if count == 0 {
		return &[]Payment{}, 0, nil
	}
	return &payments, count, nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}

func (r *Repository) findBy(sql string, params ...interface{}) ([]Payment, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return []Payment{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query(sql, params...)
	if err != nil {
		return []Payment{}, 0, err
	}

	var payments []Payment
	for rows.Next() {
		var payment Payment

		err = rows.Scan(
			&payment.UserID,
			&payment.ID,
			&payment.ProductID,
			&payment.PaymentCode,
			&payment.Price,
			&payment.Currency,
			&payment.Created,
			&payment.Updated,
		)
		if err != nil {
			return []Payment{}, 0, err
		}
		payments = append(payments, payment)
	}
	return payments, len(payments), nil
}
