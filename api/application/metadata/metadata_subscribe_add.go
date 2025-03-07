/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据添加订阅
 */
package metadata

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// AddSubscribeParams 智能元数据添加订阅参数
type AddSubscribeParams struct {
	// IP地址，通知接收地址（必填）
	//
	// 订阅成功后，SDC将向此地址推送订阅消息；
	// 12.0.0版本起支持域名，
	// 域名长度最大48个字节，即IP地址或域名入参二选一。
	Address string `json:"address"`
	// 参数通知接收端口（必填）
	//
	// 订阅成功后，SDC将向此端口推送订阅消息
	Port int `json:"port"`
	// 订阅剩余时间（必填）
	//
	// 订阅成功后，订阅服务进入倒计时，超时后将不在向目的地址推送事件，
	// 可在超时之前调用此接口刷新订阅时间，单位秒。
	//
	// 取值范围：0~86400，0表示永久有效
	TimeOut int `json:"timeOut"`
	// 是否使用https协议（必填）
	//
	// 由于元数据为敏感数据，只支持通过https协议上报。
	//
	// 此字段固定填1
	HttpsEnable int `json:"httpsEnable"`
	// 服务端开放的元数据接收HTTP/HTTPS服务URL（选填）
	MetaDataURL string `json:"metaDataURL,omitempty"`
	// 元数据上报时的认证用户名（选填）
	//
	// 用户名和密码只支持通过HTTPS加密传输，
	// 若增加、更新订阅时采用HTTP链接，则忽略digUserName和digUserPwd字段信息
	DigUserName string `json:"digUserName,omitempty"`
	// 元数据URL 登入密码（选填）
	//
	// 用户名和密码只支持通过HTTPS加密传输，
	// 若增加、更新订阅时采用HTTP链接，则忽略digUserName和digUserPwd字段信息
	DigUserPwd string `json:"digUserPwd,omitempty"`
	// 是否上报图片（选填）
	//
	// 不携带此字段或此字段设置为1时上报图片
	//
	// 取值范围：0-不上报图片，1-上报图片
	NeedUploadPic *int `json:"needUploadPic,omitempty"`
}

// AddSubscribeReplyData 智能元数据添加订阅响应数据
type AddSubscribeReplyData struct {
	// 通用响应状态
	common.ResponseStatus
	// 订阅ID。
	//
	// 订阅成功后返回，用于唯一标识此订阅，后续订阅操作（查、删、改）均基于此ID进行
	ID int `json:"ID"`
}

// AddSubscribeReply 智能元数据添加订阅响应
type AddSubscribeReply = common.Response[AddSubscribeReplyData]

// AddSubscribe 智能元数据添加订阅
//
//	@param params: 订阅参数
//	@return 订阅ID
func (p *Manager) AddSubscribe(params AddSubscribeParams) (int, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply AddSubscribeReply
	_, err := client.Post("/SDCAPI/V2.0/Metadata/Subscription").
		SetJSON(&params).
		DecodeJSON(&reply)
	if err != nil {
		return 0, err
	}

	// 检查是否订阅成功
	if reply.ResponseStatus.StatusCode != 0 {
		return 0, errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return reply.ResponseStatus.ID, nil
}
