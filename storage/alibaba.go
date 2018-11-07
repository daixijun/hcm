package storage

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AlibabaCloudOSSBackend struct {
	Client        *oss.Client
	Bucket        *oss.Bucket
	RootDirectory string
}

func NewAlibabaCloudOSSBackend(endpoint, accessKeyID, accessKeySecret, bucketName, rootDirectory string) *AlibabaCloudOSSBackend {
	if len(endpoint) == 0 {
		endpoint = "oss-cn-hangzhou.aliyuncs.com"
	}

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic("Failed to create OSS client: " + err.Error())
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic("Failed to get OSS bucket: " + err.Error())
	}

	return &AlibabaCloudOSSBackend{
		Client:        client,
		Bucket:        bucket,
		RootDirectory: rootDirectory,
	}
}

func (backend AlibabaCloudOSSBackend) IsExist(path string) bool {
	isExist, err := backend.Bucket.IsObjectExist(path)
	if err != nil {
		panic("Failed to check object is exist: " + err.Error())
	}
	return isExist
}

func (backend AlibabaCloudOSSBackend) VerifyDigest(path string, digest string) bool {
	meta, err := backend.Bucket.GetObjectDetailedMeta(path)
	if err != nil {
		panic("Failed to get object meta: " + err.Error())
	}
	meta.Del("x-oss-meta-x-oss-meta-digest")
	if meta.Get("x-oss-meta-digest") != digest {
		return false
	}

	return true
}

func (backend AlibabaCloudOSSBackend) Present(path string, data []byte, digest string) error {

	options := []oss.Option{
		oss.Meta("digest", digest),
		oss.ContentType("application/x-tar"),
	}
	return backend.Bucket.PutObject(path, bytes.NewReader(data), options...)
}
