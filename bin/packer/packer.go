package main

import (
	"fmt"
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
	if (host == "" || appid == "" || src == "" || dst == "") {
		os.Exit(1)
	}
	chunks, extraChunk := elf.SplitELF(src)
	md5, _ := common.GetFileMD5(src)
	fmt.Printf("md5: %s\n", md5)
	headerChunks := elf.HeaderChunks([]byte(appid), []byte(md5))
	var cryptoTable storage.CryptoTable
	client.FetchPubKey(host, appid, &cryptoTable)
	// fmt.Printf("pub key: %s\n", cryptoTable.PubKey)
	// headerChunks + chunks + extraChunk
	fullChunks := append(headerChunks, chunks...)
	fullChunks = append(fullChunks, extraChunk)
	//strings.ReplaceAll(cryptoTable.PubKey, "\n", "")
	encryptChunks := elf.EncryptChunks(fullChunks, []byte(cryptoTable.PubKey))

	elf.WriteChunk(dst, encryptChunks)
}
