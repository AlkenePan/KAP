package elf

import (
	"fmt"
	"os"
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
	chunks, extraChunk := SplitELF("/opt/youzu/main")
	ChunkDumper(extraChunk)
	pri, pub := crypto.GenerateKeyPair(2048)
	headerChunks := HeaderChunks([]byte("35097736-fb36-4c7e-9217-61794e9299dc"), []byte("6c6f31e624e2094dae7db53772855140"))
	// headerChunks + chunks + extraChunk
	fullChunks := append(headerChunks, chunks...)
	fullChunks = append(fullChunks, extraChunk)
	encryptChunks := EncryptChunks(fullChunks, crypto.PublicKeyToBytes(pub))
	for _, chunk := range encryptChunks[:10] {
		ChunkDumper(chunk)
	}
	ChunkDumper(fullChunks[len(fullChunks)-1])
	WriteChunk("/opt/youzu/main.encrypted", encryptChunks)
	_, _, elfBytes, _ := LoadEncryptedFile("/opt/youzu/main.encrypted", crypto.PrivateKeyToBytes(pri))
	f, err := os.Open("/opt/youzu/main")
	check(err)
	defer f.Close()
	tmpBytes := make([]byte, len(elfBytes))
	f.Read(tmpBytes)
	for i:=0;i<len(elfBytes);i++ {
		if elfBytes[i] != tmpBytes[i] {
			fmt.Println(i)
		}
	}

	//fmt.Println(len(elf))

}
