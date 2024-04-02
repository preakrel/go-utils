//go:build darwin
// +build darwin

package shell

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"
)

func GetShell(conn net.Conn) *exec.Cmd {
	_, _ = conn.Write([]byte("You're in `bash -i` shell,enter `exit` to exit\n"))
	cmd := exec.Command("/bin/bash", "-i")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

const (
	_PROC_PIDVNODEPATHINFO = 9
)

type vnodePathInfo struct {
	pvi_cdir  uint64
	pvi_cdiri uint32
	pvi_rdir  uint64
	pvi_rdiri uint32
}

func getChildrenProcesses(parentPID int) ([]int, error) {
	children := make([]int, 0)
	// Open the parent process directory
	dirPath := fmt.Sprintf("/proc/%d", parentPID)
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the vnode path information
	var pathInfo vnodePathInfo
	_, _, errno := unix.Syscall6(
		unix.SYS_PROC_INFO,
		_PROC_PIDVNODEPATHINFO,
		uintptr(parentPID),
		0,
		uintptr(unsafe.Pointer(&pathInfo)),
		0,
		0,
	)
	if errno != 0 {
		return nil, errno
	}

	// Open the vnode directory
	vnodeDirPath := fmt.Sprintf("/dev/fd/%d", pathInfo.pvi_cdiri)
	vnodeDir, err := os.Open(vnodeDirPath)
	if err != nil {
		return nil, err
	}
	defer vnodeDir.Close()

	// Read the entries in the vnode directory
	entryNames, err := vnodeDir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// Iterate through entries and parse process IDs
	for _, entryName := range entryNames {
		pid, err := strconv.Atoi(entryName)
		if err == nil {
			children = append(children, pid)
		}
	}
	return children, nil
}

// KillByParentId 杀掉进程 根据父进程id
func KillByParentId(parentProcessID int) {
	childrenIds, err := getChildrenProcesses(parentProcessID)
	if err != nil {
		return
	}
	//defer unix.Kill(parentProcessID, unix.SIGTERM)
	defer exec.Command("pkill", "-P", strconv.Itoa(parentProcessID)).Run()
	for _, childId := range childrenIds {
		//_ = unix.Kill(childId, unix.SIGTERM)
		_ = exec.Command("kill", "-9", strconv.Itoa(childId)).Run()
	}
}
