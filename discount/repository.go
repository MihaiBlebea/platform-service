package discount

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

	tableSchema := `CREATE TABLE IF NOT EXISTS discounts(
		id INT NOT NULL AUTO_INCREMENT,
		product_id INT(100),
		code VARCHAR(100) UNIQUE NOT NULL,
		percentage FLOAT(10) NOT NULL,
		active TINYINT(1) NOT NULL DEFAULT 1,
		expires TIMESTAMP NOT NULL,
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

// Add adds a Discount to the db
func (r *Repository) Add(discount *Discount) (int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare("INSERT INTO discounts (product_id, code, percentage, active, expires, created) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(
		discount.ProductID,
		discount.Code,
		discount.Percentage,
		discount.Active,
		discount.Expires,
		discount.Created,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Remove removes a discount from db
func (r *Repository) Remove(discount *Discount) error {
	client, err := r.Connection.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	stmt, err := client.Prepare("DELETE FROM discounts WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(discount.ID)
	if err != nil {
		return err
	}

	return nil
}

// FindByProductID return discounts by product
func (r *Repository) FindByProductID(productID int) (discounts []Discount, count int, err error) {
	discounts, count, err = r.findBy("SELECT * FROM discounts WHERE product_id = ?", productID)
	if err != nil {
		return discounts, 0, err
	}
	if count == 0 {
		return discounts, 0, nil
	}
	return discounts, count, nil
}

// FindByCode return discounts by discount code
func (r *Repository) FindByCode(code string) (discount Discount, count int, err error) {
	discounts, count, err := r.findBy("SELECT * FROM discounts WHERE code = ?", code)
	if err != nil {
		return discount, 0, err
	}
	if count == 0 {
		return discount, 0, nil
	}
	return discounts[0], 1, nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}

func (r *Repository) findBy(sql string, params ...interface{}) (discounts []Discount, count int, err error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return discounts, 0, err
	}
	defer client.Close()

	rows, err := client.Query(sql, params...)
	if err != nil {
		return discounts, 0, err
	}

	for rows.Next() {
		var discount Discount

		err = rows.Scan(
			&discount.ID,
			&discount.ProductID,
			&discount.Code,
			&discount.Percentage,
			&discount.Active,
			&discount.Expires,
			&discount.Created,
			&discount.Updated,
		)
		if err != nil {
			return discounts, 0, err
		}
		discounts = append(discounts, discount)
	}
	return discounts, len(discounts), nil
}
