package pkg

import (
	"io/ioutil"
	"os"
	"strings"
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
