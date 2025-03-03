/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备注册管理
 */
package devicemanage

import (
	"net/http"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// DeviceVersion 设备版本信息对象。
type DeviceVersion struct {
	Software string `json:"Software"` // 软件版本。
	Uboot    string `json:"Uboot"`    // Uboot版本。
	Kernel   string `json:"Kernel"`   // 内核版本。
	Hardware string `json:"Hardware"` // 硬件版本。
}

// ChannelInfo 视频通道信息。
type ChannelInfo struct {
	ChannelId int    `json:"ChannelId"` // 视频通道ID。
	UUID      string `json:"UUID"`      // 视频通道UUID。
	DeviceId  string `json:"DeviceId"`  // 设备ID。
}

// InitiativeRegisterParams 主动注册请求参数。
type InitiativeRegisterParams struct {
	DeviceName     string        `json:"DeviceName"`     // 设备名称。
	Manufacturer   string        `json:"Manufacturer"`   // 厂商。
	DeviceType     string        `json:"DeviceType"`     // 款型名称。
	SerialNumber   string        `json:"SerialNumber"`   // 设备序列号。
	DeviceVersion  DeviceVersion `json:"DeviceVersion"`  // 设备版本信息对象。
	IpAddr         string        `json:"IpAddr"`         // 设备IP。
	ChannelInfoArr []ChannelInfo `json:"ChannelInfoArr"` // 视频通道信息数组。
}

// InitiativeRegisterReply 设备主动注册响应参数
type InitiativeRegisterReply = common.Response[common.ResponseStatus]

// InitiativeRegister 设备主动注册（该接口通常无需外部调用）
func (p *Manager) InitiativeRegister() (*InitiativeRegisterParams, error) {
	// 获取Socket连接
	connInfo := p.connInstance.Lock()
	defer p.connInstance.Unlock()

	// 读取设备注册信息
	var params InitiativeRegisterParams
	reader := connInfo.HttpServer().Reader()
	err := reader.BindJSON(&params)
	if err != nil {
		// 构建通用响应
		res := common.NewResponseWithFailed(reader.RawRequest())
		// 响应失败结果
		err = connInfo.HttpServer().Writer().
			SetStatusCode(http.StatusOK).
			JSON(res)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	// 构建通用响应
	res := common.NewResponseWithSuccess(reader.RawRequest())
	// 响应成功结果
	err = connInfo.HttpServer().Writer().
		SetStatusCode(http.StatusOK).
		JSON(res)
	if err != nil {
		return nil, err
	}

	// OK
	return &params, nil
}
