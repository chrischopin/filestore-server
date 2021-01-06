package db

import (
	"database/sql"
	mysqldb "filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinishedToDB(filehash string, filename string, filesize int64, fileaddr string) bool {
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

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMetadataFromDB(filehash string) (*TableFile, error) {
	statement, err := mysqldb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer statement.Close()

	tableFile := TableFile{}
	err = statement.QueryRow(filehash).Scan(&tableFile.FileHash, &tableFile.FileAddr,
		&tableFile.FileName, &tableFile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tableFile, nil
}
