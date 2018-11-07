package parser

import (
	"github.com/daixijun/hcm/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

// 解析index.yaml索引文件
func ParseIndex(u string, client http.Client) (*models.IndexBody, error) {
	sep := ""
	ret := &models.IndexBody{}

	if !strings.HasSuffix(u, "/") {
		sep = "/"
	}
	indexUrl := strings.Join([]string{u, "index.yaml"}, sep)
	r, err := client.Get(indexUrl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	err = yaml.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
