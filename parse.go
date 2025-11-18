package lsblkjson_parse

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Runninginsilence1/lsblkjson_parse/internal/utils"
	"github.com/Runninginsilence1/lsblkjson_parse/model"
)

// ReadForensicDiskContextCallback 监听信号以中止检测USB设备
// 通过 callback 执行逻辑
func ReadForensicDiskContextCallback(ctx context.Context, callback func([]model.Blockdevice)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			blockDevices, _ := ReadForensicDisk()
			if len(blockDevices) > 0 {
				callback(blockDevices)
			} else {

			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func ReadForensicDiskContext(ctx context.Context) {

}

// ReadForensicDisk 检测USB设备
func ReadForensicDisk() ([]model.Blockdevice, error) {
	var blockdeviceList model.Blockdevices
	var returnDeviceList = make([]model.Blockdevice, 0)
	outContext, err := RawDeviceInfo()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(outContext, &blockdeviceList)
	if err != nil {
		err = fmt.Errorf("JSON 解码出错: %w", err)
		return returnDeviceList, err
	}
	for n, device := range blockdeviceList.Blockdevices {
		if utils.MatchMountPoint(device.Mountpoint) {
			foramtDevice := utils.DealDevice(device, n)
			returnDeviceList = append(returnDeviceList, foramtDevice)
			continue
		}
		foramtDevice := utils.DealPartition(device, n)

		if len(foramtDevice.Children) != 0 {
			returnDeviceList = append(returnDeviceList, foramtDevice)
		}
	}
	return returnDeviceList, nil
}
