package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"../meta"
	"../util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			_, _ = io.WriteString(w, "internel server error\n")
			fmt.Printf("Failed to start server,err:%s \n",err.Error())
			return
		}
		_, _ = io.WriteString(w, string(data))
	}else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file,head,err := r.FormFile("file")
		if err!=nil {
			fmt.Printf("Failed to get data,err:%s \n",err.Error())
			return
		}
		defer file.Close()//关闭文件

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/"+head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//文件流
		newFile,err := os.Create(fileMeta.Location)
		if err!=nil {
			fmt.Printf("Failed to create file,err:%s \n",err.Error())
			return
		}
		defer newFile.Close()

		//拷贝到新文件的buffer区
		fileMeta.FileSize,err = io.Copy(newFile,file)
		if err!=nil {
			fmt.Printf("Failed to save data into file,err:%s \n")
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		//流程走完成功上传 重定向
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

//处理文件上传成功信息
func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	_, _ = io.WriteString(w, "Upload success!")
}
