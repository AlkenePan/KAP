package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	MFD_CREATE  = 319
	MFD_CLOEXEC = 0x0001
)

type MemFD struct {
	*os.File
}

func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func CreateMemfd(name string) *MemFD {
	fd, _, _ := syscall.Syscall(MFD_CREATE, uintptr(unsafe.Pointer(&name)), uintptr(MFD_CLOEXEC), 0)
	return &MemFD{
		os.NewFile(fd, name),
	}
}

func (self *MemFD) Write(bytes []byte) (int, error) {
	return syscall.Write(int(self.Fd()), bytes)
}

func (self *MemFD) Path() string {
	return fmt.Sprintf("/proc/self/fd/%d", self.Fd())
}

func (self *MemFD) ExecuteWithAttributes(procAttr *syscall.ProcAttr, arguments ...string) (int, uintptr, error) {
	return syscall.StartProcess(self.Path(), append([]string{self.Name()}, arguments...), procAttr)
}

func ExecveMemfd(path string, argv []string, envv []string) error {
	data, err := ReadFile(path)
	path_list := strings.Split(path, "/")
	file_name := path_list[len(path_list)]

	if err != nil {
		return err
	}

	memfd := CreateMemfd(file_name)
	_, err = memfd.Write(data)
	if err != nil {
		return err
	}

	os_envv := os.Environ()

	for i := range envv {
		os_envv = append(os_envv, envv[i])
	}

	err = syscall.Exec(memfd.Path(), argv, os_envv)
	return err
}
