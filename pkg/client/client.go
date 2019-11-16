package client

import (
	"encoding/json"
	"net/http"
	"time"
	"youzoo/why/pkg/storage"
)
var httpClient = &http.Client{Timeout: 10 * time.Second}


func FetchPubKey(host string, appid string, target interface{}) error {
	resp, err := httpClient.Get(
		"http://" + host + "/key/pub/" + appid,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}


func FetchPriKey(host string, appid string, target interface{}) error {
	resp, err := httpClient.Get(
		"http://" + host + "/key/pri/" + appid,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}


func FetchAppInfo(host string, appid string, target interface{}) error {
	resp, err := httpClient.Get(
		"http://" + host + "/app/" + appid,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}