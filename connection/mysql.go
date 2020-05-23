package connection

import (
	"database/sql"
	"fmt"
	"os"
)

// MysqlConnection holds the mysql connection
type MysqlConnection struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// Connection interface
type Connection interface {
	Connect() (*sql.DB, error)
}

// Connect returns a mysql connection
func (mc *MysqlConnection) Connect() (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mc.Username,
		mc.Password,
		mc.Host,
		mc.Port,
		mc.Database,
	)
	client, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Mysql returns a new MysqlConnection struct
func Mysql() *MysqlConnection {
	return &MysqlConnection{
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	}
}
