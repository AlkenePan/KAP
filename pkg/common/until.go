package common

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
)

var (
	GreenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	WhiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	YellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	RedBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	BlueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	MagentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	CyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Green        = string([]byte{27, 91, 51, 50, 109})
	White        = string([]byte{27, 91, 51, 55, 109})
	Yellow       = string([]byte{27, 91, 51, 51, 109})
	Red          = string([]byte{27, 91, 51, 49, 109})
	Blue         = string([]byte{27, 91, 51, 52, 109})
	Magenta      = string([]byte{27, 91, 51, 53, 109})
	Cyan         = string([]byte{27, 91, 51, 54, 109})
	Reset        = string([]byte{27, 91, 48, 109})
	DisableColor = false
)

func GetDNSServer() string {
	f, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return ""
	}

	ori_data, err := ioutil.ReadAll(f)
	ori_str := string(ori_data)
	ori_list := strings.Split(ori_str, "\n")
	dns := strings.ReplaceAll(ori_list[0], "nameserver", "")
	return strings.TrimSpace(dns)
}

func GetBytesMD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetFileMD5(fileAbsPath string) (string, error) {
	f, err := os.Open(fileAbsPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, _ := f.Stat()
	si := fi.Size()
	fileBytes := make([]byte, si)
	f.Read(fileBytes)
	return GetBytesMD5(fileBytes), nil
}