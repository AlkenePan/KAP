package client

import (
	"fmt"
	"testing"
	"youzoo/why/pkg/app"
)

func TestFetchAppInfo(t *testing.T) {
	appinfo := app.App{}
	FetchAppInfo("localhost:5000", "679a2d46-8313-48f2-8d4b-915107b78eda", &appinfo)
	fmt.Println(appinfo)
}