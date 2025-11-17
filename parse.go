package lsblkjson_parse

import (
	"context"
	"time"

	"github.com/Runninginsilence1/lsblkjson_parse/model"
	"github.com/Runninginsilence1/lsblkjson_parse/utils"
)

// ReadStorageDevicesContext 监听信号以中止检测USB设备
func ReadStorageDevicesContext(ctx context.Context, callback func([]model.Blockdevice)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			blockDevices := utils.ReadForensicDisk()
			if len(blockDevices) > 0 {
				callback(blockDevices)
			} else {

			}
			time.Sleep(2 * time.Second)
		}
	}()
}
