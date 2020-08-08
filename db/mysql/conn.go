package mysql

import (
	"database/sql"
	"fmt"
	"os"


	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init(){
	db,_=sql.Open("mysql","root:LNMP123@tcp(54.95.127.139:33306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err:=db.Ping()
	if err!=nil {
		fmt.Printf("Failed to connect to mysql,err"+err.Error())
		os.Exit(1)
	}
}

//外部接口 返回数据库连接对象
func DBConn()*sql.DB{
	return db
}
