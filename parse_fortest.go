//go:build !linux

package lsblkjson_parse

import "os"

func RawDeviceInfo() ([]byte, error) {
	file, err := os.ReadFile("lsblk_testdata_utuntu24.json")
	if err != nil {
		return nil, err
	}
	return file, nil
}
