package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	PTRACE       = 101
	MFD_CREATE   = 319
	MFD_CLOEXEC  = 0x0001
	PTRACE_SEIZE = 0x4206
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

func DenyPtrace(pid int) (err error) {
	_, _, e := syscall.Syscall6(PTRACE, uintptr(PTRACE_SEIZE), uintptr(pid), uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if e != 0 {
		err = syscall.Errno(e)
		return err
	}
	return
}

func ExecveMemfdFromBytes(process_name string, data []byte, user_name string, ptrace bool, argv []string, envv []string) (*os.Process, error) {
	user_info, err := user.Lookup(user_name)
	if err != nil {
		return nil, err
	}

	uid, _ := strconv.Atoi(user_info.Uid)
	gid, _ := strconv.Atoi(user_info.Gid)

	memfd := CreateMemfd(process_name)
	_, err = memfd.Write(data)
	if err != nil {
		return nil, err
	}

	os_envv := os.Environ()

	for i := range envv {
		os_envv = append(os_envv, envv[i])
	}

	procAttr := &os.ProcAttr{
		Env:   os_envv,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			},
			Setsid: true,
		},
	}

	process, err := os.StartProcess(memfd.Path(), append([]string{process_name}, argv...), procAttr)

	if err != nil {
		return process, err
	}

	if !ptrace {
		err = DenyPtrace(process.Pid)
		if err != nil {
			_ = process.Kill()
			return process, err
		}
	}

	return process, err
}

func ExecveMemfdFromFile(path string, user_name string, ptrace bool, argv []string, envv []string) (*os.Process, error) {
	data, err := ReadFile(path)

	if err != nil {
		return nil, err
	}

	user_info, err := user.Lookup(user_name)
	if err != nil {
		return nil, err
	}

	uid, _ := strconv.Atoi(user_info.Uid)
	gid, _ := strconv.Atoi(user_info.Gid)

	memfd := CreateMemfd(path)
	_, err = memfd.Write(data)
	if err != nil {
		return nil, err
	}

	os_envv := os.Environ()

	for i := range envv {
		os_envv = append(os_envv, envv[i])
	}

	procAttr := &os.ProcAttr{
		Env:   os_envv,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			},
			Setsid: true,
		},
	}

	process, err := os.StartProcess(memfd.Path(), append([]string{path}, argv...), procAttr)

	if err != nil {
		return process, err
	}

	if !ptrace {
		err = DenyPtrace(process.Pid)
		if err != nil {
			_ = process.Kill()
			return process, err
		}
	}

	return process, err
}
