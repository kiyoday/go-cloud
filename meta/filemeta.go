package meta

import (
	"../db"
	"fmt"
)

//FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	//map初始化
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: 新增、更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1]=fmeta
	fmt.Printf("update sha:\n %s \n",fmeta.FileSha1)
}

// 上传meta信息入库
func UpdateFileMetaDB(fmeta FileMeta) bool{
	res := db.OnFIleUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
	return res
}

// GetFileMeta: 通过SHA1 获取文件信息对象
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

//通过SHA1 获取文件信息对象
func GetFileMetaDB(fileSha1 string) (*FileMeta, error){
	tfile,err := db.GetFileMeta(fileSha1)
	if tfile == nil ||err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}

	return &fmeta,err
}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}