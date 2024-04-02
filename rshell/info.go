package shell

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Info struct {
	mux       sync.Mutex
	KeepAlive bool
	close     bool
	pipeR     net.Conn
	pipeW     net.Conn
	// 在读取原始连接的时候往pipe写入,如果没有对这个shell做IO操作,那么会阻塞,造成shell断开后,mshell无法感知.
	// 增加一个chan 做缓冲, 避免阻塞
	buffChan chan []byte
}

func NewInfo() *Info {
	info := &Info{KeepAlive: false, buffChan: make(chan []byte, 1<<10)}
	info.pipeR, info.pipeW = net.Pipe()
	return info
}

func (i *Info) SetReadDeadline(t time.Time) error {
	return i.pipeR.SetReadDeadline(t)
}

func (i *Info) SetWriteDeadline(t time.Time) error {
	return i.pipeW.SetReadDeadline(t)
}

func (i *Info) Inject(b []byte) {
	n := len(b)
	if n == 1 {
		if b[0] == '\a' {
			i.KeepAlive = true
			return
		}
	}
	_, _ = i.pipeW.Write(b[:n])
}

func (i *Info) Monitor() {
	for {
		select {
		case b, ok := <-i.buffChan:
			if !ok {
				return
			}
			i.Inject(b)
		}
	}
}

func (i *Info) Exec(cmd string) (err error) {
	cmd = strings.Trim(cmd, "\n")
	cmd = strings.Trim(cmd, "\r")
	if i.GetClose() {
		return errors.New("close")
	}
	i.buffChan <- []byte(fmt.Sprintln(cmd))
	return
}

func (i *Info) KeyCombine(key byte) (err error) {
	i.buffChan <- []byte{key}
	return err
}

func (i *Info) GetClose() bool {
	i.mux.Lock()
	defer i.mux.Unlock()
	return i.close
}

func (i *Info) Close() (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if i.close {
		return nil
	}
	i.close = true
	close(i.buffChan)
	i.pipeW.Close()
	i.pipeR.Close()
	return
}

func (i *Info) Read(b []byte) (int, error) {
	return i.pipeR.Read(b)
}

func (i *Info) Write(b []byte) (int, error) {
	return i.pipeW.Write(b)
}

func (i *Info) ClearBuf() {
loop:
	_ = i.pipeR.SetReadDeadline(time.Now().Add(time.Second))
	_, err := i.pipeR.Read(make([]byte, 1024))
	if err == nil {
		goto loop
	}
	_ = i.pipeR.SetReadDeadline(time.Time{})
}
