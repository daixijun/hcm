package models

import "time"

type IndexBody struct {
	ApiVersion string              `yaml:"apiVersion"`
	Entries    map[string][]*Entry `yaml:"entries"`
}

type Entry struct {
	Name        string              `yaml:"name"`
	Digest      string              `yaml:"digest"`
	Urls        []string            `yaml:"urls"`
	ApiVersion  string              `yaml:"apiVersion"`
	AppVersion  string              `yaml:"appVersion"`
	Created     time.Time           `yaml:"created"`
	Description string              `yaml:"description"`
	Home        string              `yaml:"home"`
	Maintainers []map[string]string `yaml:"maintainers"`
	Sources     []string            `yaml:"sources"`
	Version     string              `yaml:"version"`
}

//type RepositoryItem struct {
//	Repo   string
//	Name   string
//	URL    string
//	Digest string
//}
