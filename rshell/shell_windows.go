//go:build windows || !linux || !darwin || !freebsd
// +build windows !linux !darwin !freebsd

package shell

import (
	"net"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"
)

func GetShell(conn net.Conn) *exec.Cmd {
	_, _ = conn.Write([]byte("You're in `cmd.exe` shell,enter `exit` to exit\n"))
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

const (
	TH32CS_SNAPPROCESS = 0x00000002
)

var (
	kernel32                 = syscall.MustLoadDLL("kernel32.dll")
	createToolhelp32Snapshot = kernel32.MustFindProc("CreateToolhelp32Snapshot")
	process32Next            = kernel32.MustFindProc("Process32NextW")
	closeHandle              = kernel32.MustFindProc("CloseHandle")
)

type PROCESSENTRY32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   uintptr
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [260]uint16
}

func createToolHelp32SnapshotCall(dwFlags, th32ProcessID uint32) (uintptr, error) {
	ret, _, err := createToolhelp32Snapshot.Call(uintptr(dwFlags), uintptr(th32ProcessID))
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func process32NextCall(hSnapshot uintptr, pe *PROCESSENTRY32) (bool, error) {
	ret, _, err := process32Next.Call(hSnapshot, uintptr(unsafe.Pointer(pe)))
	if ret == 0 {
		if err.Error() == "The handle is invalid." {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// KillByParentId 杀掉进程 根据父进程id
func KillByParentId(parentProcessID int) {
	hSnapshot, err := createToolHelp32SnapshotCall(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return
	}
	defer func() {
		_ = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(parentProcessID)).Run() //杀掉父进程
		closeHandle.Call(hSnapshot)
	}()
	var pe PROCESSENTRY32
	pe.dwSize = uint32(unsafe.Sizeof(pe))

	// Retrieve information about the first process in the snapshot
	ok, err := process32NextCall(hSnapshot, &pe)
	for ok && err == nil {
		// Check if the process has the specified parent process ID
		if pe.th32ParentProcessID == uint32(parentProcessID) {
			_ = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(int(pe.th32ProcessID))).Run()
		}
		ok, err = process32NextCall(hSnapshot, &pe)
	}
}
