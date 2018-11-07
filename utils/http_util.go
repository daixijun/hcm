package utils

import (
	"log"
	"net/http"
	"net/url"
)

func NewClient(proxy string) http.Client {
	var client http.Client

	if proxy != "" {
		proxy, err := url.Parse(proxy)
		if err != nil {
			log.Fatal(err)
		}
		trans := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		client = http.Client{
			Transport: trans,
			//Timeout:   10 * time.Second,
		}
	} else {
		//client = http.Client{Timeout: 10 * time.Second}
		client = http.Client{}
	}
	return client
}
