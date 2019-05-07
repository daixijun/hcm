package main

import (
	"github.com/daixijun/hcm/config"
	"github.com/daixijun/hcm/models"
	"github.com/daixijun/hcm/parser"
	"github.com/daixijun/hcm/storage"
	"github.com/daixijun/hcm/utils"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

// This variable is replaced in compile time
// `-ldflags "-X 'github.com/daixijun/hcm/cmd/hcm/main.Version=${VERSION}'"`
var (
    Version = "v0.0.1"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	app := cli.NewApp()
	app.Name = "hcm"
	app.Usage = "Helm Chart Mirror"
	app.Author = "Xijun Dai"
	app.Email = "daixijun1990@gmail.com"
	app.HelpName = "hcm"
	app.Version = Version
	app.Action = cliHandler
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yaml",
			Usage: "load configuration from file",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose mode",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func cliHandler(c *cli.Context) error {
	var wg sync.WaitGroup
	var backend storage.Backend

	if c.Bool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	conf := config.NewConfig(c.String("config"))
	log.Infof("Load configuration from %s\n", c.String("config"))

	client := utils.NewClient(conf.HTTPProxy)

	backend = storage.NewStorageBackend(conf)
	log.Infof("Use %s to storage.\n", conf.Storage.Backend)

	for _, repo := range conf.Repositories {
		wg.Add(1)
		go Mirror(repo, client, backend, &wg)
	}

	wg.Wait()
	return nil
}

// Mirror 同步仓库
func Mirror(repository config.Repository, client http.Client, backend storage.Backend, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		newIndexData   models.IndexBody
		chartIndexData *models.IndexBody
		chartIndexPath = path.Join(repository.Name, "index.yaml")
		chartName      string
		chartFilePath  string
		err            error
	)

	chartIndexData, err = parser.ParseIndex(repository.URL, client)
	if err != nil {
		log.Errorf("Failed to parse index %s: %s", repository.Name, err.Error())
	}

	err = copier.Copy(&newIndexData, &chartIndexData)
	if err != nil {
		log.Warnf("copy new index data failed %s: %s", repository.Name, err.Error())
	}

	for _, entryVersions := range chartIndexData.Entries {
		for index, entry := range entryVersions {
			chartName = path.Base(entry.Urls[0])
			chartFilePath = path.Join(repository.Name, chartName)

			if backend.IsExist(chartFilePath) && backend.VerifyDigest(chartFilePath, entry.Digest) {
				log.Debugf("[skip] %s \n", chartFilePath)
			} else {
				data, err := fetchChartPackage(entry.Urls[0], client)
				if err != nil {
					log.Warnf("Failed to fetch chart package %s: %s", entry.Urls[0], err.Error())
				}
				err = backend.Present(chartFilePath, data, entry.Digest)
				if err != nil {
					log.Warnf("Failed to present %s", chartFilePath)
				}
				log.Infof("[Add] %s \n", chartFilePath)
			}
			newIndexData.Entries[entry.Name][index].Urls[0] = chartName
		}
	}

	data, err := yaml.Marshal(newIndexData)
	if err != nil {
		log.Warnf("Failed to marshal index: %s", err.Error())
	}

	err = backend.Present(chartIndexPath, data, "")
	if err != nil {
		log.Warnf("Failed to present index %s: %s", chartIndexPath, err.Error())
	}
	log.Infof("[Add] %s", chartIndexPath)
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
