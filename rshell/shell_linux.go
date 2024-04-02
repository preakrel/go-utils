//go:build linux || freebsd || !windows
// +build linux freebsd !windows

package shell

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func GetShell(conn net.Conn) *exec.Cmd {
	_, _ = conn.Write([]byte("You're in `bash -i` shell,enter `exit` to exit\n"))
	cmd := exec.Command("/bin/bash", "-i")
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL, Setpgid: true}
	return cmd
}

func getChildrenProcesses(parentPID int) ([]int, error) {
	children := make([]int, 0)
	// Read entries from /proc directory
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc directory: %v", err)
	}
	// Iterate through entries in /proc
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		ppid, err := getParentProcessID(pid)
		if err != nil {
			continue
		}
		if ppid == parentPID {
			children = append(children, pid)
		}
	}
	return children, nil
}

func getParentProcessID(pid int) (int, error) {
	// Open /proc/<pid>/status file
	path := filepath.Join("/proc", strconv.Itoa(pid), "status")
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read status file for process %d: %v", pid, err)
	}
	// Find PPid in the content
	var ppid int
	_, err = fmt.Sscanf(string(content), "PPid:\t%d", &ppid)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PPid for process %d: %v", pid, err)
	}
	return ppid, nil
}

// KillByParentId 杀掉进程 根据父进程id
func KillByParentId(parentProcessID int) {
	childrenIds, err := getChildrenProcesses(parentProcessID)
	if err != nil {
		return
	}
	defer exec.Command("pkill", "-P", strconv.Itoa(parentProcessID)).Run()
	for _, childId := range childrenIds {
		_ = exec.Command("kill", "-9", strconv.Itoa(childId)).Run()
	}
}
