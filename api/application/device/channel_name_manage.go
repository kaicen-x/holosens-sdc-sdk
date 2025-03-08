/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备通道名称管理
 */
package device

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// ChannelNameInfo 通道名称信息
type ChannelNameInfo struct {
	UUID        string `json:"UUID"`        // 通道UUID
	ChannelName string `json:"channelName"` // 通道名称，1~127位字符
}

// IdQueryReply 设备通道名称查询响应
type ChannelNameQueryReply = []ChannelNameInfo

// ChannelNameQuery 设备通道名称查询
//
//	@param	uuid: 通道UUID
func (p *Manager) ChannelNameQuery(uuid string) (ChannelNameQueryReply, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply ChannelNameQueryReply
	_, err := client.Get("/SDCAPI/V1.0/CnsPaas/ChnQury/CnsChnParam").
		SetQuery("uuid", uuid).
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return reply, nil
}

// ChannelNameSetting 设备通道名称配置参数
type ChannelNameSettingParams struct {
	ChannelName string `json:"channelName"` // 通道名称，1~127位字符
}

// ChannelNameSettingReply 设备通道名称配置响应
type ChannelNameSettingReply = common.Response[common.ResponseStatus]

// ChannelNameSetting 设备通道名称配置
//
//	@param	uuid: 通道UUID
//	@param	params: 配置参数
func (p *Manager) ChannelNameSetting(uuid string, params ChannelNameSettingParams) error {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply ChannelNameSettingReply
	_, err := client.Put("/SDCAPI/V1.0/CnsPaas/ChnQury/CnsChnParam").
		SetQuery("uuid", uuid).
		SetJSON(&params).
		DecodeJSON(&reply)
	if err != nil {
		return err
	}

	// 检查状态码
	if reply.ResponseStatus.StatusCode != 0 {
		return errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return nil
}
