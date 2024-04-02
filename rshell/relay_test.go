package shell

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestShell(t *testing.T) {
	lay := NewRelay()
	lay.ClearBuf()
	lay.SetRunCmdDir("/")

	if err := lay.Exec("tree C:/"); err != nil {
		t.Errorf("执行命令失败 [%s]", err.Error())
	}

	if err := lay.Run(); err != nil {
		t.Error(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	select {
	case ers, ok := <-lay.GetErrs():
		if ok {
			t.Errorf("执行任务失败:" + ers.Error())
			return
		}
	case <-ctx.Done():

		t.Log("执行成功 监听输出")

		out := make(chan []byte, 1)
		go lay.ReadShell(out) //获取命令行输出

		for {
			select {
			case o, ok := <-out:
				if !ok {
					return
				}
				fmt.Printf("%s", o)
			case <-lay.Ctx.Done():
				return
			}
		}
	}
}
