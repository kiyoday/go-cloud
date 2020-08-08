package handler

import (
	mydb "../db/mysql"
	"encoding/json"
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
			fmt.Printf("Failed to save data into file,err:%s \n",err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetaDB(fileMeta)

		//流程走完成功上传 重定向
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

//处理文件上传成功信息
func UploadSucHandler(w http.ResponseWriter, r *http.Request){
	_, _ = io.WriteString(w, "Upload success!")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta,err := meta.GetFileMetaDB(filehash)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to get file meta,err:%s \n",err.Error())
		return
	}

	data, err := json.Marshal(fMeta) // 转换成json格式
	if err!=nil {
		fmt.Printf("failed to get file meta,err:%s \n",err.Error())
	}
	w.Write(data)

}

//文件下载接口
func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f,err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data,err := ioutil.ReadAll(f)
	if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)

}

// 更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	//meta.UpdateFileMeta(curFileMeta)
	meta.UpdateFileMetaDB(curFileMeta)

	// TODO: 更新文件表中的元信息记录

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	// 删除文件 物理删除
	os.Remove(fMeta.Location)
	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)
	// TODO: 删除表文件信息

	w.WriteHeader(http.StatusOK)
}

// 测试数据库连接
func DbTest(w http.ResponseWriter, r *http.Request) {
	err := mydb.DBConn().Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to connect database,err:%s \n",err.Error())
		return
	}else{
		data, _ := json.Marshal("ok")
		w.Write(data)
		return
	}
}