package token

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"

	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Repository struct
type Repository struct {
	Connection c.Connection
}

// Migrate generates a tokens table
func (r *Repository) Migrate() error {
	client, err := r.Connection.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	tableSchema := `CREATE TABLE IF NOT EXISTS tokens(
		id INT NOT NULL AUTO_INCREMENT,
		user_id INT(100),
		token VARCHAR(250) NOT NULL,
		hash VARCHAR(250) NOT NULL,
		active INT(1) DEFAULT 1,
		host VARCHAR(200) DEFAULT NULL,
		expires VARCHAR(50),
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

// Add adds a Token model to the db
func (r *Repository) Add(token *Token) (int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare("INSERT INTO tokens (user_id, token, hash, active, host, expires, created) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(token.UserID, token.Token, token.Hash, token.Active, token.LinkedHost, token.Expires, token.Created)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// RemoveByUser removes all tokens for user
func (r *Repository) RemoveByUser(userID int) error {
	client, err := r.Connection.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	stmt, err := client.Prepare("DELETE FROM tokens WHERE user_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

// FindByUserID returns all tokens for a user
func (r *Repository) FindByUserID(userID int) (*[]Token, int, error) {
	tokens, count, err := r.findBy("SELECT * FROM tokens WHERE user_id = ?", userID)
	if err != nil {
		return &[]Token{}, 0, err
	}
	if count == 0 {
		return &[]Token{}, 0, nil
	}
	return &tokens, count, nil
}

// FindToken returns a token by it's string value
func (r *Repository) FindToken(token string) (*Token, int, error) {

	tokens, count, err := r.findBy("SELECT * FROM tokens WHERE token = ?", token)
	if err != nil {
		return &Token{}, 0, err
	}
	if count == 0 {
		return &Token{}, 0, nil
	}
	return &tokens[0], 1, nil
}

// Update CRUD operation for Token
func (r *Repository) Update(token *Token) (int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return 0, err
	}
	defer client.Close()

	stmt, err := client.Prepare("UPDATE tokens SET user_id=?, active=?, host=?, expires=? WHERE id=?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(token.UserID, token.Active, token.LinkedHost, token.Expires, token.ID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Repo returns a repository
func Repo(conn c.Connection) *Repository {
	return &Repository{conn}
}

func (r *Repository) findBy(sql string, params ...interface{}) ([]Token, int, error) {
	client, err := r.Connection.Connect()
	if err != nil {
		return []Token{}, 0, err
	}
	defer client.Close()

	rows, err := client.Query(sql, params...)
	if err != nil {
		return []Token{}, 0, err
	}

	var tokens []Token
	for rows.Next() {
		var token Token
		err = rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.Hash,
			&token.Active,
			&token.LinkedHost,
			&token.Expires,
			&token.Created,
			&token.Updated,
		)
		if err != nil {
			return tokens, 0, err
		}
		tokens = append(tokens, token)
	}

	return tokens, len(tokens), nil
}
