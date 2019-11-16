package elf

import (
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"youzoo/why/pkg/crypto"
)

const (
	MagicNumber = "31F"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func OpenFile(file string) io.ReaderAt {
	r, err := os.Open(file)
	check(err)
	return r
}

func LoadELF(absFilePath string) *elf.File {
	f := OpenFile(absFilePath)
	err := MagicNumberCheck(f)
	check(err)
	_elf, err := elf.NewFile(f)
	check(err)
	return _elf
}

func MagicNumberCheck(f io.ReaderAt) error {
	// Read and decode ELF identifier
	var ident [16]uint8
	_, err := f.ReadAt(ident[0:], 0)
	check(err)
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		err := fmt.Errorf("Bad magic number at %d\n", ident[0:4])
		return err
	}
	return nil
}

type Chunk struct {
	Data        []byte
	EncryptFlag bool
}
type Chunks []Chunk

// Chunks, LastChunk
func SplitELF(fileAbsPath string) (Chunks, Chunk) {
	f, err := os.Open(fileAbsPath)
	check(err)
	defer f.Close()
	// file length
	fi, err := f.Stat()
	check(err)
	size := fi.Size()
	extraBytes := size % 64
	times := (size - extraBytes) / 64

	// extraBytes will append bottom
	var chunks Chunks
	for i := int64(0); i < times; i++ {
		tmpBytes := make([]byte, 64)
		_, err := f.Read(tmpBytes)
		check(err)
		chunks = append(chunks, Chunk{Data: tmpBytes, EncryptFlag: Even(i)})
	}
	tmpBytes := make([]byte, extraBytes)
	_, err = f.Read(tmpBytes)
	if err != nil {
		if err.Error() != "EOF" {
			check(err)
		}
	}
	lastChunk := Chunk{tmpBytes, false}
	return chunks, lastChunk
}

func WriteChunk(fileAbsPath string, chunks Chunks) {
	f, err := os.Create(fileAbsPath)
	check(err)
	defer f.Close()

	for _, chunk := range chunks {
		_, err := f.Write(chunk.Data)
		check(err)
	}
}

func Even(number int64) bool {
	return number%2 == 0
}
func Odd(number int64) bool {
	return !Even(number)
}

func ChunkDumper(chunk Chunk) {
	// dumper
	fmt.Printf("%d bytes:\n", len(chunk.Data))
	dumper := hex.Dumper(os.Stdout)
	defer dumper.Close()
	dumper.Write(chunk.Data)
}

func EncryptChunk(chunk Chunk, pubKeyBytes []byte) Chunk {
	pubKey := crypto.BytesToPublicKey(pubKeyBytes)
	encryptData := crypto.EncryptWithPublicKey(chunk.Data, pubKey)
	chunk.Data = encryptData
	return chunk
}

func DecryptChunk(chunk Chunk, priKeyBytes []byte) Chunk {
	priKey := crypto.BytesToPrivateKey(priKeyBytes)
	decryptData := crypto.DecryptWithPrivateKey(chunk.Data, priKey)
	chunk.Data = decryptData
	return chunk

}

func EncryptChunks(chunks Chunks, pubKeyBytes []byte) Chunks {
	var encryptChunks Chunks
	for _, chunk := range chunks {
		if chunk.EncryptFlag {
			encryptedChunk := EncryptChunk(chunk, pubKeyBytes)
			encryptChunks = append(encryptChunks, encryptedChunk)

		} else {
			encryptChunks = append(encryptChunks, chunk)
		}
	}
	return encryptChunks
}

func HeaderChunks(appid []byte, hash []byte) Chunks {
	magicNumber := []byte(MagicNumber)
	magicChunk := Chunk{magicNumber, false}
	appidChunk := Chunk{appid, false}
	hashChunk := Chunk{hash, true}
	return Chunks{magicChunk, appidChunk, hashChunk}
}

func IntToBytes(number uint32) []byte {
	bs := make([]byte, 5)
	binary.LittleEndian.PutUint32(bs, number)
	return bs
}
func BytesToInt(bs []byte) uint32 {
	number := binary.LittleEndian.Uint32(bs)
	return number
}

/*
3 bytes:
00000000  33 31 46                                          |31F|
5 bytes:
00000000  06 00 00 00 00                                    |.....|
6 bytes:
00000000  31 32 33 31 32 33                                 |123123|
256 bytes:
00000000  55 07 5b 00 46 50 1a f7  2f f7 37 54 1e bf 15 8e  |U.[.FP../.7T....|
00000010  62 1d a8 9c 7a 96 89 ef  de 83 3c 1d 91 9f 51 05  |b...z.....<...Q.|
00000020  d0 31 fa 08 f6 ff a5 03  9b ab d2 20 f2 79 74 d1  |.1......... .yt.|
00000030  bd 36 95 13 16 f3 ee 89  eb fb 7e ac 88 89 d1 d3  |.6........~.....|
*/
func LoadEncryptedFileHeader(fileAbsPath string) (string, []byte) {
	f, err := os.Open(fileAbsPath)
	check(err)
	defer f.Close()
	magicNumber := make([]byte, 3)
	_, err = f.Read(magicNumber)
	check(err)
	// appid
	appid := make([]byte, 36)
	_, err = f.Read(appid)
	check(err)
	// hash
	hash := make([]byte, 256)
	_, err = f.Read(hash)
	check(err)
	return string(appid), hash
}

func LoadEncryptedFile(fileAbsPath string, priKey []byte) (string, []byte, []byte, error) {
	f, err := os.Open(fileAbsPath)
	check(err)
	defer f.Close()
	magicNumber := make([]byte, 3)
	_, err = f.Read(magicNumber)
	check(err)
	// appid
	appid := make([]byte, 36)
	_, err = f.Read(appid)
	check(err)
	// hash
	hash := make([]byte, 256)
	_, err = f.Read(hash)
	check(err)

	fi, _ := f.Stat()
	size := fi.Size()
	restFileLength := size - int64(3+36+256)
	var chunks Chunks
	next := 256
	for {
		var tmpBytes []byte
		if restFileLength < 64 {
			tmpBytes = make([]byte, restFileLength)
			f.Read(tmpBytes)
			chunk := Chunk{tmpBytes, false}
			chunks = append(chunks, chunk)
			break
		} else {
			tmpBytes = make([]byte, next)
		}
		_, err := f.Read(tmpBytes)
		if err != nil {
			check(err)
		}
		if next == 256 {
			chunk := Chunk{tmpBytes, true}
			if restFileLength > 256 {
				restFileLength -= 256
			}
			next = 64
			chunks = append(chunks, chunk)

		} else {
			chunk := Chunk{tmpBytes, false}
			if restFileLength > 64 {
				restFileLength -= 64
			}
			next = 256
			chunks = append(chunks, chunk)

		}

	}
	//ChunkDumper(chunks[len(chunks)-1])
	var plainELF []byte
	for _, chunk := range chunks {
		if chunk.EncryptFlag {
			encryptedChunk := DecryptChunk(chunk, priKey)
			plainELF = append(plainELF, encryptedChunk.Data...)

		} else {
			plainELF = append(plainELF, chunk.Data...)

		}
	}
	return string(appid), hash, plainELF, nil
}
