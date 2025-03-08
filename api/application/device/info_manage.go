/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备信息管理
 */
package device

import (
	"strconv"
)

// OsMemInfo 操作系统内存详细统计信息
type OsMemInfo struct {
	TotalMem int `json:"totalMem"` // 可用总内存（KB）
	FreeMem  int `json:"freeMem"`  // 空闲内存（KB）
	BuffMem  int `json:"buffMem"`  // Buff内存（KB）
	CacheMem int `json:"cacheMem"` // Cache内存（KB）
}

// MemoryObject 内存信息对象
type MemoryObject struct {
	DdrMem    int       `json:"ddrMem"`    // DDR内存总容量（单位：GB）
	OsMem     int       `json:"osMem"`     // 操作系统分配的内存大小（单位：MB）
	OsMemInfo OsMemInfo `json:"osMemInfo"` // 操作系统内存详细统计信息
}

// MMZ 内存对象
type MMZInfo struct {
	TotalMem int `json:"totalMem"` // 可用总内存（KB）
	FreeMem  int `json:"freeMem"`  // 空闲内存（KB）
	UsedMem  int `json:"usedMem"`  // 占用内存（KB）
}

// Flash 信息对象
type FlashInfo struct {
	TotalFlashSize int `json:"totalFlashSize"` // Flash总大小（KB）
	FreeFlashSize  int `json:"freeFlashSize"`  // Flash空闲大小（KB）
}

// NetRateInfo 网络统计信息对象
type NetRateInfo struct {
	RecvRate int `json:"recvRate"` // 网络接收速率，单位bit/s
	SendRate int `json:"sendRate"` // 网络发送速率，单位bit/s
}

// BaseInfoQueryReply 设备基础信息查询响应
type BaseInfoQueryReply struct {
	PlatformType     string       `json:"platformType"`     // 设备平台类型（如Linux/iOS等）
	BarCode          string       `json:"barCode"`          // 设备唯一条形码标识
	BomCode          string       `json:"bomCode"`          // 设备BOM编码（物料清单编号）
	DrvCode          string       `json:"drvCode"`          // 设备驱动标识码
	DevType          string       `json:"devType"`          // 设备款型名称（如DS-2CD2020G0-I）
	DzoomRatio       float64      `json:"dzoomRatio"`       // 数字变倍倍率（浮点数表示）
	CpuOccupyRate    float64      `json:"cpuOccupyRate"`    // CPU占用率百分比（浮点数）
	FlashVersion     string       `json:"flashVersion"`     // Flash版本信息后缀（格式：V.X，最大8位字符）
	FullDeviceType   string       `json:"fullDeviceType"`   // 完整设备型号（格式：外部型号@内部型号，最大128位）
	PackageLimitSize int          `json:"packageLimitSize"` // 升级包限制大小（前端校验用，整型）
	MemInfo          MemoryObject `json:"memInfo"`          // 内存信息对象（包含DDR/OS内存详情）
	MmzInfo          MMZInfo      `json:"mmzInfo"`          // MMZ专用内存信息对象
	FlashInfo        FlashInfo    `json:"flashInfo"`        // Flash存储信息对象
	NetRateInfo      NetRateInfo  `json:"netRateInfo"`      // 网络统计信息对象（参考表2-24）
	Manufacturer     string       `json:"manufacturer"`     // 厂商名称（SDC 11.1.0版本新增）
	SoftVersion      string       `json:"softVersion"`      // 软件包版本信息
	KernelVersion    string       `json:"kernelVersion"`    // 内核版本（SDC 11.1.0版本新增）
	ESN              string       `json:"ESN"`              // 设备序列号（ESN号，SDC 11.1.0版本新增）
	Eth0Mac          string       `json:"eth0Mac"`          // 网卡1 MAC地址（SDC 11.1.0版本新增）
	Eth1Mac          string       `json:"eth1Mac"`          // 网卡2 MAC地址（SDC 11.1.0版本新增）
	UbootVersion     string       `json:"ubootVersion"`     // Uboot版本信息（SDC 11.1.0版本新增）
	HardVersion      string       `json:"hardVersion"`      // 硬件版本信息（SDC 11.1.0版本新增）
}

// BaseInfoQuery 设备基础信息查询
//
//	@param channelID: 通道ID，针对复眼款型可用，普通款型无需传入此参数，或传入101。取值范围：101 - 定点信息，102- 复眼全景路信息。
func (p *Manager) BaseInfoQuery(channelID int) (*BaseInfoQueryReply, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply BaseInfoQueryReply
	_, err := client.Get("/SDCAPI/V1.0/MiscIaas/System").
		SetQuery("ChannelID", strconv.Itoa(channelID)).
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
