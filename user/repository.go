package user

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

	tableSchema := `CREATE TABLE IF NOT EXISTS users(
		id int NOT NULL AUTO_INCREMENT,
		name VARCHAR(250),
		email VARCHAR(100) UNIQUE,
		password VARCHAR(100),
		jwt_token VARCHAR(255) DEFAULT NULL,
		active INT(1) DEFAULT 1,
		consent INT(1) DEFAULT 0,
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
func (r *Repository) Add(user *User) (int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare("INSERT INTO users (name, email, password, jwt_token, active, consent, created) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.JWT, user.Active, user.Consent, user.Created)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// FindByID returns a user by id
func (r *Repository) FindByID(userID int) (*User, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return &User{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query("SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return &User{}, 0, err
	}

	var users []User
	for rows.Next() {
		var user User

		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.JWT, &user.Active, &user.Consent, &user.Created, &user.Updated)
		if err != nil {
			return &User{}, 0, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return &User{}, 0, nil
	}
	return &users[0], len(users), nil
}

// FindByEmail returns a user by id
func (r *Repository) FindByEmail(email string) (*User, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return &User{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return &User{}, 0, err
	}

	var users []User
	for rows.Next() {
		var user User

		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.JWT, &user.Active, &user.Consent, &user.Created, &user.Updated)
		if err != nil {
			return &User{}, 0, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return &User{}, 0, nil
	}
	return &users[0], len(users), nil
}

// FindByJWT returns a user based on a jwt token
func (r *Repository) FindByJWT(token string) (*User, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return &User{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query("SELECT * FROM users WHERE jwt_token = ?", token)
	if err != nil {
		return &User{}, 0, err
	}

	var users []User
	for rows.Next() {
		var user User

		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.JWT, &user.Active, &user.Consent, &user.Created, &user.Updated)
		if err != nil {
			return &User{}, 0, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return &User{}, 0, nil
	}
	return &users[0], len(users), nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}
