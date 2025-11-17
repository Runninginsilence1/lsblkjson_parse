package model

type Blockdevices struct {
	Blockdevices []Blockdevice `json:"blockdevices"`
}

// Device 统一的设备结构体，用于块设备和分区
// 兼容不同版本的 lsblk 输出（某些字段在不同版本中类型可能为 string 或 number）
type Device struct {
	Name       string   `json:"name"`
	Pname      string   `json:"pname"` // 父设备名称（仅分区使用）
	Path       string   `json:"path"`
	Fsavail    any      `json:"fsavail"`     // 可用空间，可能是 string 或 number
	FsavailRaw uint64   `json:"fsavail_raw"` // 转换后的原始数值
	Fssize     any      `json:"fssize"`      // 总大小，可能是 string 或 number
	FssizeRaw  uint64   `json:"fssize_raw"`  // 转换后的原始数值
	Fsused     any      `json:"fsused"`      // 已使用空间，可能是 string 或 number
	Fstype     string   `json:"fstype"`
	Fsuse      string   `json:"fsuse%"`
	Mountpoint string   `json:"mountpoint"`
	Label      string   `json:"label"`
	Children   []Device `json:"children"` // 子设备（分区）
}

// Blockdevice 块设备（为了兼容性保留，实际使用 Device）
type Blockdevice = Device

// Partition 分区（为了兼容性保留，实际使用 Device）
type Partition = Device

type FileWithTime struct {
	Name    string
	ModTime string
}
