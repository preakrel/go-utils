package shell

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
	"util/str"
)

type Relay struct {
	pipeR     net.Conn //命令行响应输出
	pipeW     net.Conn //命令行响应输入
	Info      *Info
	mu        sync.Mutex
	cmd       *exec.Cmd
	cmdDir    string
	Ctx       context.Context
	CancelFn  context.CancelFunc
	closeChan chan bool
	errsChan  chan error
	isClose   bool
}

func NewRelay() *Relay {
	lay := &Relay{Info: NewInfo(), closeChan: make(chan bool), errsChan: make(chan error, 1)}
	lay.Ctx, lay.CancelFn = context.WithCancel(context.Background())
	lay.pipeR, lay.pipeW = net.Pipe()
	return lay
}

func (lay *Relay) GetErrs() <-chan error {
	return lay.errsChan
}

func (lay *Relay) GetCloseChan() <-chan bool {
	return lay.closeChan
}

func (lay *Relay) Run() (err error) {
	defer func() {
		if r := recover(); r != nil && !lay.GetIsClose() {
			lay.errsChan <- errors.New(fmt.Sprint(r))
		}
	}()
	go lay.Info.Monitor()
	go func() {
		defer lay.Close()
		if err := lay.RunShell(lay.Info.pipeR); err != nil && !lay.GetIsClose() {
			lay.errsChan <- err
		}
	}()
	return
}

func (lay *Relay) SetRunCmdDir(dir string) {
	lay.cmdDir = dir
}

func (lay *Relay) RunShell(reader io.ReadWriteCloser) (err error) {
	lay.cmd = GetShell(lay.pipeW)
	if lay.cmdDir != "" {
		lay.cmd.Dir = lay.cmdDir
	}
	lay.cmd.Stdout = lay.pipeW
	lay.cmd.Stderr = lay.pipeW
	lay.cmd.Stdin = bufio.NewReader(reader)
	return lay.cmd.Run()
}

func (lay *Relay) SetReadDeadline(t time.Time) error {
	return lay.pipeR.SetReadDeadline(t)
}

func (lay *Relay) SetWriteDeadline(t time.Time) error {
	return lay.pipeW.SetReadDeadline(t)
}

func (lay *Relay) GetIsClose() bool {
	lay.mu.Lock()
	defer lay.mu.Unlock()
	return lay.isClose
}

func (lay *Relay) Close() (err error) {
	lay.mu.Lock()
	defer lay.mu.Unlock()
	defer lay.CancelFn()
	if !lay.isClose {
		lay.isClose = true
		close(lay.closeChan)
		close(lay.errsChan)

	}
	if lay.cmd != nil && lay.cmd.Process != nil {
		pid := lay.cmd.Process.Pid
		_ = lay.cmd.Process.Signal(syscall.SIGTERM) //windows 平台不支持 故通过下面协程杀进程
		go KillByParentId(pid)
		lay.cmd = nil
	}
	_ = lay.pipeR.Close()
	_ = lay.pipeW.Close()
	return lay.Info.Close()
}

func (lay *Relay) Reset() (err error) {
loop:
	_ = lay.pipeR.SetReadDeadline(time.Now().Add(time.Second))
	_, err = lay.pipeR.Read(make([]byte, 1024))
	if err == nil {
		goto loop
	}
	_ = lay.pipeR.SetReadDeadline(time.Time{})
	lay.pipeR, lay.pipeW = net.Pipe()
	return
}

func (lay *Relay) ClearBuf() {
loop:
	_ = lay.pipeR.SetReadDeadline(time.Now().Add(time.Second))
	_, err := lay.pipeR.Read(make([]byte, 1024))
	if err == nil {
		goto loop
	}
	_ = lay.pipeR.SetReadDeadline(time.Time{})
}

func (lay *Relay) Exec(cmd string) (err error) {
	return lay.Info.Exec(cmd)
}

func (lay *Relay) KeyCombine(key byte) (err error) {
	return lay.Info.KeyCombine(key)
}

func (lay *Relay) Read(b []byte) (int, error) {
	return lay.pipeR.Read(b)
}

func (lay *Relay) Write(b []byte) (int, error) {
	lay.mu.Lock()
	defer lay.mu.Unlock()
	return lay.pipeW.Write(b)
}

func (lay *Relay) ReadShell(out chan<- []byte) {
	readByte := make([]byte, 4096) //避免读取截断，TODO 长度大于此字节也会存在问题，后期优化
	defer func() { readByte = nil }()
	var msgBuf bytes.Buffer
	defer msgBuf.Reset()
	var isWait bool
	go func() { //大于4096字节 拼接字节 最后错误等待 超时，则直接输出
		timer := time.NewTicker(100 * time.Millisecond)
		defer timer.Stop()
		for {
			select {
			case <-lay.Ctx.Done():
				return
			case <-timer.C:
				if isWait && msgBuf.Len() > 0 {
					out <- msgBuf.Bytes()
					msgBuf.Reset()
				}
			}
		}
	}()
	for {

		select {
		case <-lay.Ctx.Done():
			return
		default:
			n, err := lay.Read(readByte)
			if err != nil {
				return
			}
			msgBuf.Write(readByte[:n])
			// UTF-8 长度为 2-4 字节
			// 因为读出来的数据流可能为被截断的UTF-8 校验是否正确utf-8 不正确则继续追加
			convertBytes := str.CharsetConvertToUTF8(msgBuf.Bytes(), "gbk")
			if r, size := utf8.DecodeRune(convertBytes); r == utf8.RuneError && size == 0 || !utf8.FullRune(convertBytes) { //不是正确的utf-8
				isWait = true
				continue
			}
			if !utf8.Valid(convertBytes) {
				isWait = true
				continue
			}
			isWait = false
			out <- convertBytes
			msgBuf.Reset()
		}
	}
}
