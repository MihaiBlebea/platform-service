package product

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

	tableSchema := `CREATE TABLE IF NOT EXISTS products(
		id int NOT NULL AUTO_INCREMENT,
		code VARCHAR(250) UNIQUE,
		name VARCHAR(100) UNIQUE,
		price FLOAT(6),
		currency VARCHAR(4),
		active INT(1) DEFAULT 1,
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

// FindByCode returns a user by id
func (r *Repository) FindByCode(code string) (*Product, int, error) {
	products, count, err := r.findBy("SELECT * FROM products WHERE code = ?", code)
	if err != nil {
		return &Product{}, 0, err
	}
	if count == 0 {
		return &Product{}, 0, nil
	}
	return &products[0], len(products), nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}

func (r *Repository) findBy(sql string, params ...interface{}) ([]Product, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return []Product{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query(sql, params...)
	if err != nil {
		return []Product{}, 0, err
	}

	var products []Product
	for rows.Next() {
		var product Product

		err = rows.Scan(&product.ID, &product.Code, &product.Name, &product.Price, &product.Currency, &product.Active, &product.Created, &product.Updated)
		if err != nil {
			return []Product{}, 0, err
		}
		products = append(products, product)
	}
	return products, len(products), nil
}
