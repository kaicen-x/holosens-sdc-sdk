/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标库删除
 */
package recognize

import (
	"errors"
	"net/url"

	"github.com/kaicen-x/holosens-sdc-sdk/api/common"
)

// TargetLibDeleteParam 目标库删除参数构造器
type TargetLibDeleteParam func(url.Values)

// TargetLibDeleteWithName 目标库删除参数：目标库名称
//
//	待删除的目标库名称，不填表示删除所有目标库
func TargetLibDeleteWithName(val string) TargetLibDeleteParam {
	return func(values url.Values) {
		values.Set("libName", val)
	}
}

// TargetLibDeleteReply 目标库删除响应
type TargetLibDeleteReply = common.Response[common.ResponseStatus]

// TargetLibDelete 目标库删除
//
//	@param	params：目标库删除参数（TargetLibDeleteWithName：待删除的目标库名称，不填表示删除所有目标库）
//	@return 异常信息
func (p *Manager) TargetLibDelete(params ...TargetLibDeleteParam) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetLibDeleteReply
	req := client.Delete("/SDCAPI/V2.0/FaceApp/FaceRecog/FaceLibs/Libs").
		SetContentType("application/x-www-form-urlencoded")
	for _, param := range params {
		param(req.GetQuery())
	}
	_, err := req.DecodeJSON(&reply)
	if err != nil {
		return err
	}

	// 检查是否删除成功
	if reply.ResponseStatus.StatusCode != 0 {
		return errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return nil
}
