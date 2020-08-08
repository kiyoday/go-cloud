package db

import (
	mydb "../db/mysql"
	"database/sql"
	"fmt"
)

// 数据库保存元信息 数据库接口
func OnFIleUploadFinished(filehash string,filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {//没有新的记录
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileName sql.NullString
	FileHash string
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//查询文件元信息  数据库接口
func GetFileMeta(filehash string)(*TableFile, error){
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

