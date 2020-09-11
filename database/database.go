package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DbConn *sql.DB

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", username, password, host)
}
func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", dsn())
	if err != nil {
		log.Fatal(err)
	}

	DbConn.SetMaxOpenConns(3)
	DbConn.SetMaxIdleConns(3)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
