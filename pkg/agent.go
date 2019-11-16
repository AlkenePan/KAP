package main

import (
	"flag"
	"fmt"
	"youzoo/why/pkg/agent"
	"youzoo/why/pkg/client"
	"youzoo/why/pkg/common"
	"youzoo/why/pkg/elf"
	"youzoo/why/pkg/storage"
)

func main() {
	var path string
	var host string
	var cryptoTable storage.CryptoTable

	flag.StringVar(&path, "p", "", "ELF Path")
	flag.StringVar(&host, "h", "", "Verify Server Host")

	if path == "" {
		fmt.Println(common.GreenBg, "[!] ERROR: ELF Path is Null", common.Reset)
		return
	}

	if host == "" {
		fmt.Println(common.GreenBg, "[!] ERROR: Verify Server Host is Null", common.Reset)
		return
	}

	appid, hash := elf.LoadEncryptedFileHeader(path)
	err := client.FetchPriKey(host, appid, &cryptoTable)
	if err != nil {
		fmt.Println(common.GreenBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	_, _, ori_elf_data, err := elf.LoadEncryptedFile(path, []byte(cryptoTable.PriKey))
	if err != nil {
		fmt.Println(common.GreenBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	file_current_md5 := common.GetBytesMD5(ori_elf_data)

	md5 := ""

	verify_app_request_data := agent.VerifyAPPRequest{}
	verify_app_request_data.DNS = common.GetDNSServer()
	verify_app_request_data.APPID = appid
	verify_app_request_data.MD5 = md5

}
