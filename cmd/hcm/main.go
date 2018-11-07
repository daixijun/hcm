package main

import (
	"github.com/daixijun/hcm/config"
	"github.com/daixijun/hcm/models"
	"github.com/daixijun/hcm/parser"
	"github.com/daixijun/hcm/storage"
	"github.com/daixijun/hcm/utils"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Name = "hcm"
	app.Version = "v0.0.1"
	app.Action = cliHandler

	err := app.Run(os.Args)
	if err != nil {
		panic(err.Error())
	}
}

func cliHandler(c *cli.Context) error {
	var wg sync.WaitGroup
	var backend storage.Backend

	conf := config.NewConfig("config.yaml")
	client := utils.NewClient(conf.HTTPProxy)

	backend = storage.NewStorageBackend(conf)

	for _, repo := range conf.Repositories {
		wg.Add(1)
		go Mirror(repo, client, backend, &wg)
	}

	wg.Wait()
	return nil
}

func Mirror(repository config.Repository, client http.Client, backend storage.Backend, wg *sync.WaitGroup) {
	defer wg.Done()
	var newIndexData models.IndexBody

	IndexData, err := parser.ParseIndex(repository.URL, client)
	if err != nil {
		panic("Failed to parse index " + repository.Name + " : " + err.Error())
	}
	err = copier.Copy(&newIndexData, &IndexData)

	if err != nil {
		panic("copy new index data failed: " + err.Error())
	}

	for _, entryVersions := range IndexData.Entries {
		for index, entry := range entryVersions {
			//fmt.Println(repository.Name, entry.Name, entry.Version, entry.Digest, entry.Urls[0])

			fileName := path.Base(entry.Urls[0])
			filePath := path.Join(repository.Name, fileName)
			if backend.IsExist(filePath) && backend.VerifyDigest(filePath, entry.Digest) {
				color.Cyan("[skip] %s \n", filePath)
			} else {
				data, err := fetchChartPackage(entry.Urls[0], client)
				if err != nil {
					panic("Failed to fetch chart package " + entry.Urls[0] + " : " + err.Error())
				}
				err = backend.Present(filePath, data, entry.Digest)
				if err != nil {
					panic("Failed to present " + filePath)
				}
				color.Green("[Add] %s \n", filePath)
			}
			newIndexData.Entries[entry.Name][index].Urls[0] = fileName
		}
	}

	data, err := yaml.Marshal(newIndexData)
	if err != nil {
		panic("Failed to marshal index: " + err.Error())
	}
	err = backend.Present(path.Join(repository.Name, "index.yaml"), data, "")
	if err != nil {
		panic("Failed to present index: " + err.Error())
	}
}

func fetchChartPackage(u string, client http.Client) ([]byte, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
