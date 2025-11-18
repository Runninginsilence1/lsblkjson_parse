package utils

import (
	"strconv"
	"strings"

	"github.com/Runninginsilence1/lsblkjson_parse/model"
	"github.com/spf13/cast"
)

// ConvertToUint64 将 any 类型（可能是 string 或 number）转换为 uint64
// 兼容不同版本的 lsblk 输出
func ConvertToUint64(value any) uint64 {
	if value == nil {
		return 0
	}
	return cast.ToUint64(value)
}

// ProcessDevice 递归处理设备及其子设备，转换字段类型并格式化
func ProcessDevice(device *model.Device) {
	// 转换并保存原始数值
	device.FsavailRaw = ConvertToUint64(device.Fsavail)
	device.FssizeRaw = ConvertToUint64(device.Fssize)
	fsusedRaw := ConvertToUint64(device.Fsused)

	// 格式化为可读字符串
	device.Fsavail = FormatFileSize(device.FsavailRaw)
	device.Fssize = FormatFileSize(device.FssizeRaw)
	device.Fsused = FormatFileSize(fsusedRaw)

	// 递归处理子设备
	for i := range device.Children {
		ProcessDevice(&device.Children[i])
	}
}

// 挂载
func MatchMountPoint(point string) bool {
	if strings.HasPrefix(point, "/media/") {
		return true
	}

	return false
}

// 解析属性
func DealDevice(device model.Blockdevice, n int) model.Blockdevice {
	var resDevice model.Blockdevice

	_ = DeepCopy(device, &resDevice)

	resDevice.Name = "磁盘" + strconv.Itoa(n)

	// 使用 cast 转换，兼容不同版本的 lsblk
	partitionFsavail := ConvertToUint64(resDevice.Fsavail)
	partitionFssize := ConvertToUint64(resDevice.Fssize)
	partitionFsused := ConvertToUint64(resDevice.Fsused)

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

// CollectMountedPartitions 递归收集所有已挂载的分区
func CollectMountedPartitions(devices []model.Device, parentName string, startIndex int) ([]model.Partition, uint64, uint64, uint64) {
	var partitions []model.Partition
	var totalFssize uint64
	var totalFsused uint64
	var totalFsavail uint64

	for index, partition := range devices {
		// 检查当前设备是否挂载
		if MatchMountPoint(partition.Mountpoint) {
			// 使用 cast 转换，兼容不同版本的 lsblk
			partitionFssize := ConvertToUint64(partition.Fssize)
			partitionFsused := ConvertToUint64(partition.Fsused)
			partitionFsavail := ConvertToUint64(partition.Fsavail)

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
			childPartitions, childFssize, childFsused, childFsavail := CollectMountedPartitions(
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
func DealPartition(device model.Blockdevice, n int) model.Blockdevice {
	var resDevice model.Blockdevice
	resDevice.Name = "磁盘" + strconv.Itoa(n)
	resDevice.Path = device.Path

	// 递归收集所有已挂载的分区
	partitions, totalFssize, totalFsused, totalFsavail := CollectMountedPartitions(
		device.Children, resDevice.Name, 0)

	resDevice.Children = partitions
	resDevice.Fssize = FormatFileSize(totalFssize)
	resDevice.Fsused = FormatFileSize(totalFsused)
	resDevice.Fsavail = FormatFileSize(totalFsavail)

	return resDevice
}
