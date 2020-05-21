package connection

import (
	"database/sql"
	"fmt"
)

const (
	username = "admin"
	password = "pass"
	host     = "127.0.0.1"
	port     = "3307"
	database = "platform"
)

// MysqlConnection holds the mysql connection
type MysqlConnection struct {
	username string
	password string
	host     string
	port     string
	database string
}

// Connection interface
type Connection interface {
	Connect() (*sql.DB, error)
}

// Connect returns a mysql connection
func (mc *MysqlConnection) Connect() (*sql.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mc.username, mc.password, mc.host, mc.port, mc.database)
	client, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Mysql returns a new MysqlConnection struct
func Mysql() *MysqlConnection {
	return &MysqlConnection{username, password, host, port, database}
}
