//go:build linux

package lsblkjson_parse

import "github.com/Runninginsilence1/lsblkjson_parse/internal/utils"

// 从lsblk中解析原始信息
func RawDeviceInfo() ([]byte, error) {
	outContext, err := utils.RunCommandWithOutput(0, "",
		"lsblk -f -J -e 7,11,3,2 -o NAME,PATH,FSAVAIL,FSSIZE,FSTYPE,FSUSED,FSUSE%,MOUNTPOINT,LABEL -b")

	return outContext, err
}
