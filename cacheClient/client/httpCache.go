package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	client *http.Client
	server string
}

func (c *httpClient) Get(key string) string {
	resp, e := c.client.Get(fmt.Sprintf("http://%s/%s", c.server, key))
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}
	return string(b)
}

func (c *httpClient) Set(key, value string) {
	req, e := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://%s/%s", c.server, key),
		strings.NewReader(value))
	if e != nil {
		log.Println(key)
		panic(e)
	}
	resp, e := c.client.Do(req)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func NewHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, server}
}
