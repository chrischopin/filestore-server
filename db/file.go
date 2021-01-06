package db

import (
	mysqldb "filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string)bool {
	statement, err := mysqldb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`," +
			"`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	defer statement.Close()

	ret, err := statement.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err {
		if rowsAffected <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before", filehash)
			return false
		}
		return true
	}
	return false
}