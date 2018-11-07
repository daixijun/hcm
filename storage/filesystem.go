package storage

import (
	"github.com/daixijun/hcm/utils"
	"io/ioutil"
	"os"
	"path"
)

type FileSystemBackend struct {
	RootDirectory string
}

func NewFileSystemBackend(rootDirectory string) *FileSystemBackend {
	return &FileSystemBackend{RootDirectory: rootDirectory}
}

func (backend FileSystemBackend) IsExist(p string) bool {
	filePath := path.Join(backend.RootDirectory, p)
	_, err := os.Stat(filePath)
	if err != nil {
		panic("Failed to check file is exists: " + err.Error())
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (backend FileSystemBackend) VerifyDigest(p string, digest string) bool {
	filePath := path.Join(backend.RootDirectory, p)
	data, err := ioutil.ReadFile(filePath)
	if err != err {
		panic("Failed to read file: " + err.Error())
	}
	sum := utils.NewSHA256Digest(data)
	if sum != digest {
		return false
	}
	return true
}

func (backend FileSystemBackend) Present(p string, data []byte, digest string) error {
	filePath := path.Join(backend.RootDirectory, p)
	err := ioutil.WriteFile(filePath, data, 0644)
	return err
}
