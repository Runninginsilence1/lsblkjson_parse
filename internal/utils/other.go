package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
