package main

import (
	"flag"
	"fmt"
	"strings"
	"youzoo/why/pkg/agent"
	"youzoo/why/pkg/app"
	"youzoo/why/pkg/client"
	"youzoo/why/pkg/common"
	"youzoo/why/pkg/elf"
	"youzoo/why/pkg/storage"
)

func getMsgFromFileWatcher(appid string, host string, msgChan chan agent.FileWatcherMsg) {
	for i := range msgChan {
		alert_data := storage.AlertTable{
			Appid:       appid,
			Level:       i.Level,
			Type:        "file watcher",
			Info:        i.Msg,
			PostContact: "",
		}

		err := client.NewAlert(host, appid, &alert_data)
		if err != nil {
			fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		}
	}
}

func main() {
	var path string
	var host string
	var argv []string
	var cryptoTable storage.CryptoTable

	flag.StringVar(&path, "path", "", "ELF Path")
	flag.StringVar(&host, "host", "", "Verify Server Host")

	flag.Parse()

	if path == "" {
		fmt.Println(common.RedBg, "[!] ERROR: ELF Path is Null", common.Reset)
		return
	}

	if host == "" {
		fmt.Println(common.RedBg, "[!] ERROR: Verify Server Host is Null", common.Reset)
		return
	}

	appid, hash := elf.LoadEncryptedFileHeader(path)
	err := client.FetchPriKey(host, appid, &cryptoTable)
	if err != nil {
		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	_, _, ori_elf_data, err := elf.LoadEncryptedFile(path, []byte(cryptoTable.PriKey))
	if err != nil {
		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	file_current_md5 := common.GetBytesMD5(ori_elf_data)

	tempChunk := elf.Chunk{hash, true}
	decryptChunk := elf.DecryptChunk(tempChunk, []byte(cryptoTable.PriKey))
	if string(decryptChunk.Data) != file_current_md5 {
		fmt.Println(common.RedBg, "[!] ERROR: ELF File MD5 Verify Error", common.Reset)
		return
	}

	appinfo := app.App{}
	err = client.FetchAppInfo(host, appid, &appinfo)
	if err != nil {
		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	if appinfo.DNS != "" {
		if appinfo.DNS != common.GetDNSServer() {
			fmt.Println(common.RedBg, "[!] ERROR: DNS Server Error", common.Reset)
			return
		}
	}

	if appinfo.ExecInfo.Argv != "" {
		argv = strings.Split(appinfo.ExecInfo.Argv, ";")
	}
	envv := strings.Split(appinfo.ExecInfo.Envv, ";")

	process, err := agent.ExecveMemfdFromBytes(path, ori_elf_data, appinfo.ExecInfo.UserName, appinfo.ExecInfo.Ptrace, argv, envv)

	if err != nil {
		if process != nil {
			_ = process.Kill()
		}

		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)

		alert_data := storage.AlertTable{
			Appid:       appid,
			Level:       "danger",
			Type:        "runtime",
			Info:        err.Error(),
			PostContact: "",
		}

		err = client.NewAlert(host, appid, &alert_data)
		if err != nil {
			fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		}

		return
	}

	alert_data := storage.AlertTable{
		Appid:       appid,
		Level:       "success",
		Type:        "runtime",
		Info:        "",
		PostContact: "",
	}

	err = client.NewAlert(host, appid, &alert_data)
	if err != nil {
		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
	}

	filrWatcherMsgChan := agent.GetNewFileWatcherMsgChan()
	go agent.AddNewWatcher(path, filrWatcherMsgChan)
	go getMsgFromFileWatcher(appid, host, filrWatcherMsgChan)

	_, err = process.Wait()

	if err != nil {
		fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)

		alert_data := storage.AlertTable{
			Appid:       appid,
			Level:       "danger",
			Type:        "runtime(wait)",
			Info:        err.Error(),
			PostContact: "",
		}

		err = client.NewAlert(host, appid, &alert_data)
		if err != nil {
			fmt.Println(common.RedBg, "[!] ERROR: "+err.Error(), common.Reset)
		}

		return
	}
}
