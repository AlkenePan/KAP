package main

import (
	"os"
	"youzoo/why/pkg/client"
	"youzoo/why/pkg/common"
	"youzoo/why/pkg/elf"
	"youzoo/why/pkg/storage"
)

func main() {
	if len(os.Args) <= 1 {
		os.Exit(1)
	}
	host := os.Args[1]
	appid := os.Args[2]
	src := os.Args[3]
	dst := os.Args[4]
	if host == "" || appid == "" || src == "" || dst == "" {
		os.Exit(1)
	}
	common.Info("AppId: " + appid)
	common.Info("reading file: " + src)
	chunks, extraChunk := elf.SplitELF(src)
	md5, _ := common.GetFileMD5(src)
	headerChunks := elf.HeaderChunks([]byte(appid), []byte(md5))
	var cryptoTable storage.CryptoTable
	common.Info("Connecting: " + host)
	_ = client.FetchPubKey(host, appid, &cryptoTable)
	fullChunks := append(headerChunks, chunks...)
	fullChunks = append(fullChunks, extraChunk)
	encryptChunks := elf.EncryptChunks(fullChunks, []byte(cryptoTable.PubKey))
	common.Info("Encrypt Done")
	elf.WriteChunk(dst, encryptChunks)
	common.Info("Out: " + dst)
}
