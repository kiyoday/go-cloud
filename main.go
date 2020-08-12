package main

import (
	"net/http"
	"./handler"
	"fmt"
)

func main(){
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	http.HandleFunc("/file/test",handler.DbTest)

	http.HandleFunc("/user/signup",handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.UserInfoHandler)


	err := http.ListenAndServe(":8090",nil)
	if err !=nil {
		fmt.Printf("Failed to start server,err:%s",err.Error())
	}
}