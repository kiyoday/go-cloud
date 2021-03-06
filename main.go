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
	err := http.ListenAndServe(":8090",nil)
	if err !=nil {
		fmt.Printf("Failed to start server,err:%s",err.Error())
	}
}