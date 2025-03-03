/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备ID管理
 */
package devicemanage

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// IdInfo ID信息
type IdInfo struct {
	UUID     string `json:"UUID"`     // 通道UUID
	DeviceID string `json:"deviceID"` // 设备ID，最大长度64字符
}

// IdQueryReply 设备ID查询响应
type IdQueryReply struct {
	IDs []IdInfo `json:"IDs"` // 设备ID对象数组，每个对象对应一个通道，每个通道对应一个设备ID
}

// IdQuery 设备ID查询
func (p *Manager) IdQuery() (*IdQueryReply, error) {
	// 获取Socket连接
	connInfo := p.connInstance.Lock()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply IdQueryReply
	_, err := connInfo.HttpClient().
		Get("/SDCAPI/V1.0/Rest/DeviceID").
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}

// IdSettingParams 设备ID配置参数
type IdSettingParams struct {
	IDs []IdInfo `json:"IDs"` // 设备ID对象数组，每个对象对应一个通道，每个通道对应一个设备ID
}

// IdSettingReply 设备ID配置响应
type IdSettingReply = common.Response[common.ResponseStatus]

// IdSetting 设备ID配置
//
//	@param	params: 配置参数
func (p *Manager) IdSetting(params IdSettingParams) error {
	// 获取Socket连接
	connInfo := p.connInstance.Lock()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply IdSettingReply
	_, err := connInfo.HttpClient().
		Put("/SDCAPI/V1.0/Rest/DeviceID").
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
