package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Runninginsilence1/lsblkjson_parse/model"
	"github.com/spf13/cast"
)

// convertToUint64 将 any 类型（可能是 string 或 number）转换为 uint64
// 兼容不同版本的 lsblk 输出
func convertToUint64(value any) uint64 {
	if value == nil {
		return 0
	}
	return cast.ToUint64(value)
}

// processDevice 递归处理设备及其子设备，转换字段类型并格式化
func processDevice(device *model.Device) {
	// 转换并保存原始数值
	device.FsavailRaw = convertToUint64(device.Fsavail)
	device.FssizeRaw = convertToUint64(device.Fssize)
	fsusedRaw := convertToUint64(device.Fsused)

	// 格式化为可读字符串
	device.Fsavail = FormatFileSize(device.FsavailRaw)
	device.Fssize = FormatFileSize(device.FssizeRaw)
	device.Fsused = FormatFileSize(fsusedRaw)

	// 递归处理子设备
	for i := range device.Children {
		processDevice(&device.Children[i])
	}
}

// ReadForensicDisk 检测USB设备
func ReadForensicDisk() []model.Blockdevice {
	var blockdeviceList model.Blockdevices
	var returnDeviceList = make([]model.Blockdevice, 0)
	outContext, _ := RunCommandWithOutput(0, "",
		"bash", "-c", "lsblk -f -J -e 7,11,3,2 -o NAME,PATH,FSAVAIL,FSSIZE,FSTYPE,FSUSED,FSUSE%,MOUNTPOINT,LABEL -b")
	err := json.Unmarshal(outContext, &blockdeviceList)
	if err != nil {
		fmt.Println("JSON 解码出错:", err)
		return returnDeviceList
	}
	for n, device := range blockdeviceList.Blockdevices {
		if matchMountPoint(device.Mountpoint) {
			foramtDevice := dealDevice(device, n)
			returnDeviceList = append(returnDeviceList, foramtDevice)
			continue
		}
		foramtDevice := dealPartition(device, n)

		if len(foramtDevice.Children) != 0 {
			returnDeviceList = append(returnDeviceList, foramtDevice)
		}
	}
	return returnDeviceList
}

// 挂载
func matchMountPoint(point string) bool {
	if strings.HasPrefix(point, "/media/") {
		return true
	}

	return false
}

// 解析属性
func dealDevice(device model.Blockdevice, n int) model.Blockdevice {
	var resDevice model.Blockdevice

	_ = DeepCopy(device, &resDevice)

	resDevice.Name = "磁盘" + strconv.Itoa(n)

	// 使用 cast 转换，兼容不同版本的 lsblk
	partitionFsavail := convertToUint64(resDevice.Fsavail)
	partitionFssize := convertToUint64(resDevice.Fssize)
	partitionFsused := convertToUint64(resDevice.Fsused)

	partition := model.Partition{
		Name:       "分区0",
		Pname:      resDevice.Name,
		Path:       resDevice.Path,
		Fsavail:    FormatFileSize(partitionFsavail),
		FsavailRaw: partitionFsavail,
		Fssize:     FormatFileSize(partitionFssize),
		FssizeRaw:  partitionFssize,
		Fsused:     FormatFileSize(partitionFsused),
		Fstype:     resDevice.Fstype,
		Fsuse:      resDevice.Fsuse,
		Mountpoint: resDevice.Mountpoint,
		Label:      resDevice.Label,
	}
	resDevice.Children = []model.Partition{partition}
	resDevice.Fssize = FormatFileSize(partitionFssize)
	resDevice.Fsused = FormatFileSize(partitionFsused)
	resDevice.Fsavail = FormatFileSize(partitionFsavail)
	return resDevice
}

// collectMountedPartitions 递归收集所有已挂载的分区
func collectMountedPartitions(devices []model.Device, parentName string, startIndex int) ([]model.Partition, uint64, uint64, uint64) {
	var partitions []model.Partition
	var totalFssize uint64
	var totalFsused uint64
	var totalFsavail uint64

	for index, partition := range devices {
		// 检查当前设备是否挂载
		if matchMountPoint(partition.Mountpoint) {
			// 使用 cast 转换，兼容不同版本的 lsblk
			partitionFssize := convertToUint64(partition.Fssize)
			partitionFsused := convertToUint64(partition.Fsused)
			partitionFsavail := convertToUint64(partition.Fsavail)

			totalFssize += partitionFssize
			totalFsused += partitionFsused
			totalFsavail += partitionFsavail

			// 设置分区信息
			partition.Fssize = FormatFileSize(partitionFssize)
			partition.FssizeRaw = partitionFssize
			partition.Fsused = FormatFileSize(partitionFsused)
			partition.Fsavail = FormatFileSize(partitionFsavail)
			partition.FsavailRaw = partitionFsavail
			partition.Pname = parentName

			if partition.Label == "" {
				partition.Name = "分区" + strconv.Itoa(startIndex+index)
			} else {
				partition.Name = partition.Label
			}

			partitions = append(partitions, partition)
		}

		// 递归处理子分区
		if len(partition.Children) > 0 {
			childPartitions, childFssize, childFsused, childFsavail := collectMountedPartitions(
				partition.Children, parentName, startIndex+index+1)
			partitions = append(partitions, childPartitions...)
			totalFssize += childFssize
			totalFsused += childFsused
			totalFsavail += childFsavail
		}
	}

	return partitions, totalFssize, totalFsused, totalFsavail
}

// 解析属性
func dealPartition(device model.Blockdevice, n int) model.Blockdevice {
	var resDevice model.Blockdevice
	resDevice.Name = "磁盘" + strconv.Itoa(n)
	resDevice.Path = device.Path

	// 递归收集所有已挂载的分区
	partitions, totalFssize, totalFsused, totalFsavail := collectMountedPartitions(
		device.Children, resDevice.Name, 0)

	resDevice.Children = partitions
	resDevice.Fssize = FormatFileSize(totalFssize)
	resDevice.Fsused = FormatFileSize(totalFsused)
	resDevice.Fsavail = FormatFileSize(totalFsavail)

	return resDevice
}
