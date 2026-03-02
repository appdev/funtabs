package storage

import (
	"funtabs-server/config"
	"log"
	"mime/multipart"
)

// Storage 定义文件存储接口，本地和 S3 均实现此接口
type Storage interface {
	// Save 保存上传的文件，dir 为子目录（如 "avatars"），返回可访问的 URL
	Save(file multipart.File, header *multipart.FileHeader, dir string) (string, error)
	// Delete 删除文件，path 为 Save 返回的 URL
	Delete(urlPath string) error
}

// S 是全局存储实例
var S Storage

// Init 根据配置初始化存储后端
func Init() {
	switch config.Cfg.Storage.Type {
	case "s3":
		S = newS3Storage()
		log.Println("存储后端: S3/对象存储")
	default:
		S = newLocalStorage()
		log.Println("存储后端: 本地文件系统")
	}
}
