/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据订阅删除
 */
package metadata

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// SubscribeDeleteParam 订阅删除参数构建器
type SubscribeDeleteParam func(url.Values)

// SubscribeDeleteWithID 订阅删除参数：订阅ID
//
//	选填，不填则删除全部订阅
func SubscribeDeleteWithID(val int) SubscribeDeleteParam {
	return func(values url.Values) {
		values.Add("ID", strconv.Itoa(val))
	}
}

// SubscribeDeleteReply 智能元数据订阅删除响应
type SubscribeDeleteReply = common.Response[common.ResponseStatus]

// SubscribeDelete 订阅删除
//
//	@param params: 订阅删除参数
//	@return 异常信息
func (p *Manager) SubscribeDelete(params ...SubscribeDeleteParam) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发起请求
	var reply SubscribeDeleteReply
	req := client.Delete("/SDCAPI/V2.0/Metadata/Subscription").
		SetContentType("application/x-www-form-urlencoded")
	for _, param := range params {
		param(req.GetQuery())
	}
	_, err := req.DecodeJSON(&reply)
	if err != nil {
		return err
	}

	// 是否删除成功
	if reply.ResponseStatus.StatusCode != 0 {
		return errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return nil
}
