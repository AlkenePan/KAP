package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
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

func NewAlert(host string, appid string, target interface{}) error {
	jsonStr, err := json.Marshal(target)
	if err != nil {
		panic(err)
	}
	resp, err := httpClient.Post(
		"http://"+host+"/alert/new",
		"application/json",
		bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
