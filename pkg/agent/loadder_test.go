package agent

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	process, err := ExecveMemfd("/root/gohack/test", "nobody", false, []string{}, []string{})
	if err != nil {
		fmt.Println("ERROR!", err)
	}
	fmt.Println(process.Pid)
	_, _ = process.Wait()
}
