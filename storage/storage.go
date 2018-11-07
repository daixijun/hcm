package storage

import (
	"github.com/daixijun/hcm/config"
)

type Backend interface {
	IsExist(path string) bool                              // 判断文件是否存在
	VerifyDigest(path string, digest string) bool          // 检验文件 digest 值是否一致
	Present(path string, data []byte, digest string) error // 保存文件
}

func NewStorageBackend(cfg config.Config) Backend {
	var backend Backend

	switch cfg.Storage.Backend {
	case "oss":
		backend = NewAlibabaCloudOSSBackend(
			cfg.Storage.AlibabaCloudOSS.Endpoint,
			cfg.Storage.AlibabaCloudOSS.AccessKeyID,
			cfg.Storage.AlibabaCloudOSS.AccessKeySecret,
			cfg.Storage.AlibabaCloudOSS.BucketName,
			cfg.Storage.AlibabaCloudOSS.RootDirectory,
		)
	case "filesystem":
		backend = NewFileSystemBackend(cfg.Storage.FileSystem.RootDirectory)
	default:
		panic("Unsupported storage backend:" + cfg.Storage.Backend)
	}
	return backend
}
