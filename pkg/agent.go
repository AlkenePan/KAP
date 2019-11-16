package main

import (
	"flag"
	"fmt"
	"youzoo/why/pkg/agent"
	"youzoo/why/pkg/common"
)

func main() {
	var path string
	var host string

	flag.StringVar(&path, "p", "", "ELF Path")
	flag.StringVar(&host, "h", "", "Verify Server Host")

	if path == "" {
		fmt.Println(common.GreenBg, "[!] ERROR: ELF Path is Null", common.Reset)
	}

	if host == "" {
		fmt.Println(common.GreenBg, "[!] ERROR: Verify Server Host is Null", common.Reset)
	}

	appid := ""
	md5 := ""

	verify_app_request_data := agent.VerifyAPPRequest{}
	verify_app_request_data.DNS = common.GetDNSServer()
	verify_app_request_data.APPID = appid
	verify_app_request_data.MD5 = md5

}
