package db

import (
	mysqldb "filestore-server/db/mysql"
	"fmt"
)

func UserSignUp(username, password string) bool {
	statement, err := mysqldb.DBConn().Prepare("insert ignore into tbl_user " +
		"(`user_name`, `user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}
	defer statement.Close()

	ret, err := statement.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert user to db, err: " + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
