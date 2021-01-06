package db

import (
	mysqldb "filestore-server/db/mysql"
	"fmt"
)

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

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

func UserSignIn(username, encpwd string) bool {
	statement, err := mysqldb.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := statement.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("Username not found: " + username)
		return false
	}

	pRows := mysqldb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	stmt, err := mysqldb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mysqldb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
