package user

import (
	"errors"

	c "github.com/MihaiBlebea/purpletree/platform/connection"

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
		confirm_code VARCHAR(8) DEFAULT "",
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

	stmt, err := client.Prepare(`
		INSERT INTO users 
			(name, email, password, jwt_token, active, consent, confirm_code, created) 
		VALUES 
			(?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(
		user.Name,
		user.Email,
		user.Password,
		user.JWT,
		user.Active,
		user.Consent,
		user.ConfirmCode,
		user.Created,
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

// Remove removes a user from the database
func (r *Repository) Remove(user *User) error {
	client, err := r.Connection.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	stmt, err := client.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update CRUD operation for User
func (r *Repository) Update(user *User) (int, error) {
	if user.ID == 0 {
		return 0, errors.New("Could not update user as it doesn't have an id")
	}

	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare(`
		UPDATE 
			users 
		SET 
			name = ?, email = ?, password = ?, jwt_token = ?, active = ?, consent = ?, confirm_code = ? 
		WHERE 
			id=?`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(
		user.Name,
		user.Email,
		user.Password,
		user.JWT,
		user.Active,
		user.Consent,
		user.ConfirmCode,
		user.ID,
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

// FindByID returns a user by id
func (r *Repository) FindByID(userID int) (user *User, count int, err error) {
	users, count, err := r.findBy("SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return user, count, err
	}
	if count == 0 {
		return user, count, nil
	}
	return &users[0], 1, nil
}

// FindByEmail returns a user by id
func (r *Repository) FindByEmail(email string) (user *User, count int, err error) {
	users, count, err := r.findBy("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return user, count, err
	}
	if count == 0 {
		return user, count, nil
	}
	return &users[0], 1, nil
}

// FindByJWT returns a user based on a jwt token
func (r *Repository) FindByJWT(token string) (user *User, count int, err error) {
	users, count, err := r.findBy("SELECT * FROM users WHERE jwt_token = ?", token)
	if err != nil {
		return user, count, err
	}
	if count == 0 {
		return user, count, nil
	}
	return &users[0], 1, nil
}

// FindByConfirmCode returns a user base on confirm code
func (r *Repository) FindByConfirmCode(confirmCode string) (user *User, count int, err error) {
	users, count, err := r.findBy("SELECT * FROM users WHERE confirm_code = ?", confirmCode)
	if err != nil {
		return user, count, err
	}
	if count == 0 {
		return user, count, nil
	}
	return &users[0], 1, nil
}

func (r *Repository) findBy(sql string, params ...interface{}) (users []User, count int, err error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return users, count, err
	}
	defer client.Close()

	rows, err := client.Query(sql, params...)
	if err != nil {
		return users, count, err
	}

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.JWT,
			&user.Active,
			&user.Consent,
			&user.ConfirmCode,
			&user.Created,
			&user.Updated,
		)
		if err != nil {
			return users, count, err
		}
		users = append(users, user)
	}

	return users, len(users), nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}
