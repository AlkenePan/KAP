package agent

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	err := ExecveMemfd("/root/gohack/test", []string{}, []string{})
	fmt.Println(err)
}
