package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// RunCommandWithOutput 通用执行命令
func RunCommandWithOutput(timeout uint64, dir string, name string, args ...string) ([]byte, error) {
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, name, args...)
		if dir != "" {
			cmd.Dir = dir
		}
		return cmd.CombinedOutput()
	} else {
		cmd := exec.Command(name, args...)
		if dir != "" {
			cmd.Dir = dir
		}
		return cmd.CombinedOutput()
	}
}

// DeepCopy 深拷贝
func DeepCopy(src, dest interface{}) error {
	// 使用JSON序列化和反序列化实现深拷贝
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(srcBytes, dest)
	if err != nil {
		return err
	}

	return nil
}

// FormatFileSize 计算大小
func FormatFileSize(size uint64) string {
	// 定义文件大小单位
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	// 处理文件大小为0的情况
	if size == 0 {
		return "0 B"
	}

	// 计算文件大小所在单位的索引
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	// 格式化文件大小
	return fmt.Sprintf("%d %s", size, units[unitIndex])
}

// CVLCScript VLCScript 创建普通用户vlc执行vlc命令

// RunCmdContext 得到输出并返回错误, 或者上下文被取消.
// Linux 上需要专门设置进程组来关闭
// 适用于 Cmd 还需要开启各种各样的子进程的情况
// 尽量不要使用 sigkill, 所以使用 sigterm
// 并且如果是被上下文取消, 返回的 err 将会是 context.Canceled
func RunCmdContext(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	var (
		b             bytes.Buffer
		err           error
		ctxCancelFlag bool
	)
	defer func() {
		if ctxCancelFlag {
			err = context.Canceled
		}
	}()
	//开辟新的线程组（Linux 平台特有的属性）
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, //使得 Shell 进程开辟新的 PGID, 即 Shell 进程的 PID, 它后面创建的所有子进程都属于该进程组
	}
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	var finish = make(chan struct{})
	defer close(finish)
	go func() {
		select {
		case <-ctx.Done(): //超时/被 cancel 结束
			//kill -(-PGID) 杀死整个进程组
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			ctxCancelFlag = true
		case <-finish: //正常结束
		}
	}()
	//wait 等待 goroutine 执行完，然后释放 FD 资源
	//这个时候再 kill 掉 shell 进程就不会再等待了，会直接返回
	if err = cmd.Wait(); err != nil {
		return nil, err
	}
	return b.Bytes(), err
}
