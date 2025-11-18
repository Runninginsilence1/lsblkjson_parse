//go:build !linux

package lsblkjson_parse

import "os"

// 这里我不懂咋注入，使用一个极端的方案：手动设置环境变量，通过
func RawDeviceInfo() ([]byte, error) {
	file, err := os.ReadFile("lsblk_testdata_utuntu24.json")
	if err != nil {
		return nil, err
	}
	return file, nil
}
