package main

import (
	"os"
	"youzoo/why/pkg/client"
	"youzoo/why/pkg/common"
	"youzoo/why/pkg/elf"
	"youzoo/why/pkg/storage"
)

func main() {
	host := os.Args[1]
	appid := os.Args[2]
	src := os.Args[3]
	dst := os.Args[4]
	chunks, extraChunk := elf.SplitELF(src)
	md5, _ := common.GetFileMD5(src)
	headerChunks := elf.HeaderChunks([]byte(appid), []byte(md5))
	var cryptoTable storage.CryptoTable
	client.FetchPubKey(host, appid, &cryptoTable)
	// headerChunks + chunks + extraChunk
	fullChunks := append(headerChunks, chunks...)
	fullChunks = append(fullChunks, extraChunk)
	encryptChunks := elf.EncryptChunks(fullChunks, []byte(cryptoTable.PubKey))
	for _, chunk := range encryptChunks[:10] {
		elf.ChunkDumper(chunk)
	}
	elf.WriteChunk(dst, encryptChunks)
}
