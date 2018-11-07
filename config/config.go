package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// 配置文件
type Config struct {
	HTTPProxy    string       `yaml:"http_proxy"`
	Repositories []Repository `yaml:"repositories"`
	Storage      storage      `yaml:"storage"`
}

// 存储
type storage struct {
	Backend         string          `yaml:"backend"`
	AlibabaCloudOSS alibabaCloudOSS `yaml:"oss"`
	FileSystem      fileSystem      `yaml:"filesystem"`
}

// oss 存储
type alibabaCloudOSS struct {
	Endpoint        string `yaml:"endpoint"`
	BucketName      string `yaml:"bucketname"`
	RootDirectory   string `yaml:"rootdirectory"`
	AccessKeyID     string `yaml:"accesskeyid"`
	AccessKeySecret string `yaml:"accesskeysecret"`
}

// 本地目录存储
type fileSystem struct {
	RootDirectory string `yaml:"rootdirectory"`
}

// 需要同步的仓库
type Repository struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func NewConfig(configFile string) Config {
	var cfg Config

	if configFile == "" {
		configFile = "cfg.yaml"
	}

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件错误: %s", err.Error())
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("解析配置文件错误， 请确认配置是否为YAML格式: %s", err.Error())
	}

	if len(cfg.Repositories) == 0 {
		log.Fatalln("`repositories` 配置不能为空")
	}

	return cfg
}
