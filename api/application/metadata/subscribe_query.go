/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据订阅查询
 */
package metadata

import (
	"net/url"
	"strconv"
)

// SubscribeQueryParam 订阅查询参数构建器
type SubscribeQueryParam func(url.Values)

// SubscribeQueryWithID 订阅查询参数：订阅ID
//
//	选填，不填则返回全部订阅
func SubscribeQueryWithID(val int) SubscribeQueryParam {
	return func(values url.Values) {
		values.Add("ID", strconv.Itoa(val))
	}
}

// SubscribeQueryReply 订阅查询响应
type SubscribeQueryReply struct {
	// 订阅信息列表
	Subscriptions []SubscribeInfo `json:"subscriptions"`
}

// SubscribeQuery 订阅查询
//
//	@param params: 订阅查询参数
//	@return 订阅查询结果
//	@return 异常信息
func (p *Manager) SubscribeQuery(params ...SubscribeQueryParam) (*SubscribeQueryReply, error) {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发起请求
	var reply SubscribeQueryReply
	req := client.Get("/SDCAPI/V2.0/Metadata/Subscription").
		SetContentType("application/x-www-form-urlencoded")
	for _, param := range params {
		param(req.GetQuery())
	}
	_, err := req.DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
