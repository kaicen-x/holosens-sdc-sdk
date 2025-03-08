/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据修改订阅
 */
package metadata

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// SubscribeChangeParams 智能元数据订阅修改参数
type SubscribeChangeParams struct {
	// 订阅信息
	SubscribeInfo
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

// SubscribeChangeReply 智能元数据订阅修改响应
type SubscribeChangeReply = common.Response[common.ResponseStatus]

// SubscribeChange 智能元数据订阅修改
//
//	@param params: 订阅参数
//	@return 订阅ID
func (p *Manager) SubscribeChange(params SubscribeChangeParams) error {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply SubscribeChangeReply
	_, err := client.Put("/SDCAPI/V2.0/Metadata/Subscription").
		SetJSON(&params).
		DecodeJSON(&reply)
	if err != nil {
		return err
	}

	// 检查是否修改成功
	if reply.ResponseStatus.StatusCode != 0 {
		return errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return nil
}
