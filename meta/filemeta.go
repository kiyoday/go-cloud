package meta

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

//UpdateFileMeta: 新增、更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1]=fmeta
}
// GetFileMeta: 通过SHA1 获取文件信息对象
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}