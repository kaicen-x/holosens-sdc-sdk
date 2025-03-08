/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据订阅添加
 */
package metadata

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// SubscribeAddParams 智能元数据订阅添加参数
type SubscribeAddParams struct {
	// 订阅基础信息
	SubscribeBaseInfo
	// 元数据上报时的认证用户名（选填）
	//
	// 用户名和密码只支持通过HTTPS加密传输，
	// 若增加、更新订阅时采用HTTP链接，则忽略digUserName和digUserPwd字段信息
	DigUserName *string `json:"digUserName,omitempty"`
	// 元数据URL 登入密码（选填）
	//
	// 用户名和密码只支持通过HTTPS加密传输，
	// 若增加、更新订阅时采用HTTP链接，则忽略digUserName和digUserPwd字段信息
	DigUserPwd *string `json:"digUserPwd,omitempty"`
}

// SubscribeAddReplyData 智能元数据订阅添加响应数据
type SubscribeAddReplyData struct {
	// 通用响应状态
	common.ResponseStatus
	// 订阅ID。
	//
	// 订阅成功后返回，用于唯一标识此订阅，后续订阅操作（查、删、改）均基于此ID进行
	ID int `json:"ID"`
}

// SubscribeAddReply 智能元数据添加订阅响应
type SubscribeAddReply = common.Response[SubscribeAddReplyData]

// SubscribeAdd 智能元数据订阅添加
//
//	@param params: 订阅添加参数
//	@return 订阅ID
func (p *Manager) SubscribeAdd(params SubscribeAddParams) (int, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply SubscribeAddReply
	_, err := client.Post("/SDCAPI/V2.0/Metadata/Subscription").
		SetJSON(&params).
		DecodeJSON(&reply)
	if err != nil {
		return 0, err
	}

	// 检查是否添加成功
	if reply.ResponseStatus.StatusCode != 0 {
		return 0, errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return reply.ResponseStatus.ID, nil
}
