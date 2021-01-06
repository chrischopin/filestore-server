package mysql

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql, err: %s\n", err.Error())
		os.Exit(1)
	}
}

// return db connection instance
func DBConn() *sql.DB {
	return db
}
