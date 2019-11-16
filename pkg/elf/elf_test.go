package elf

import (
	"fmt"
	"testing"
	"youzoo/why/pkg/crypto"
)

func TestELF(t *testing.T) {
	//_ := LoadELF("/opt/youzu/smith")
	chunks, extraChunk := SplitELF("/opt/youzu/ls")
	for _, chunk64 := range chunks {
		ChunkDumper(chunk64)
	}
	ChunkDumper(extraChunk)
	chunks = append(chunks, extraChunk)
	WriteChunk("/opt/youzu/ls_bak", chunks)
}

func TestEncrypt(t *testing.T) {
	chunks, extraChunk := SplitELF("/opt/youzu/ls")
	priv, pub := crypto.GenerateKeyPair(2048)
	encrptyChunk := EncryptChunk(extraChunk, crypto.PublicKeyToBytes(pub))
	decryptChunk := DecryptChunk(encrptyChunk, crypto.PrivateKeyToBytes(priv))
	fmt.Println(extraChunk)
	fmt.Println(encrptyChunk)
	fmt.Println(decryptChunk)
	for _, chunk64 := range chunks {
		encrptyChunk := EncryptChunk(chunk64, crypto.PublicKeyToBytes(pub))
		fmt.Println(len(chunk64.Data), len(encrptyChunk.Data))
	}
}

func TestEncryptChunks(t *testing.T) {
	chunks, extraChunk := SplitELF("/opt/youzu/ls")
	_, pub := crypto.GenerateKeyPair(2048)
	headerChunks := HeaderChunks([]byte("123123"), []byte("asdfasfsafasdf"))
	// headerChunks + chunks + extraChunk
	fullChunks := append(headerChunks, chunks...)
	fullChunks = append(fullChunks, extraChunk)
	encryptChunks := EncryptChunks(fullChunks, crypto.PublicKeyToBytes(pub))
	for _, chunk := range encryptChunks[:10] {
		ChunkDumper(chunk)
	}
	WriteChunk("/opt/youzu/ls.encrypted", encryptChunks)
	//elf, _ := LoadEncryptedFile("/opt/youzu/ls.encrypted", crypto.PrivateKeyToBytes(pri))
	//fmt.Println(len(elf))

}

