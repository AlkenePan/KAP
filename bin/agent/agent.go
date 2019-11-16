package agent

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

	tempChunk := elf.Chunk{hash, true}
	decryptChunk := elf.DecryptChunk(tempChunk, []byte(cryptoTable.PriKey))
	if string(decryptChunk.Data) != file_current_md5 {
		fmt.Println(common.GreenBg, "[!] ERROR: ELF File MD5 Verify Error", common.Reset)
		return
	}

	appinfo := app.App{}
	err = client.FetchAppInfo(host, appid, &appinfo)
	if err != nil {
		fmt.Println(common.GreenBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	if appinfo.DNS != common.GetDNSServer() {
		fmt.Println(common.GreenBg, "[!] ERROR: DNS Server Error", common.Reset)
		return
	}

	argv := strings.Split(appinfo.ExecInfo.Argv, ";")
	envv := strings.Split(appinfo.ExecInfo.Envv, ";")

	process, err := agent.ExecveMemfdFromBytes(path, ori_elf_data, appinfo.ExecInfo.UserName, appinfo.ExecInfo.Ptrace, argv, envv)

	if err != nil {
		if process != nil {
			_ = process.Kill()
		}

		fmt.Println(common.GreenBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}

	_, err = process.Wait()

	if err != nil {
		fmt.Println(common.GreenBg, "[!] ERROR: "+err.Error(), common.Reset)
		return
	}
}
