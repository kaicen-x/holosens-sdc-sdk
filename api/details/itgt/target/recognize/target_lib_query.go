/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标库查询
 */
package recognize

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// TargetLibQueryReplyData 目标库查询响应数据
type TargetLibQueryReplyData struct {
	FaceLibNum int                    `json:"FaceLibNum"`    // 目标库总数
	TargetLibs []TargetLibDetailsInfo `json:"FaceListsArry"` // 目标库列表
}

// TargetLibQueryReply 目标库查询响应
type TargetLibQueryReply struct {
	// 通用响应状态
	common.ResponseStatus
	// 目标库查询响应数据
	TargetLibQueryReplyData
}

// TargetLibQuery 目标库查询
//
//	@return 目标库查询响应数据
//	@return 异常信息
func (p *Manager) TargetLibQuery() (*TargetLibQueryReplyData, error) {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetLibQueryReply
	_, err := client.Get("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/Libs").
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// 检查是否获取成功
	if reply.StatusCode != 0 {
		return nil, errors.New(reply.StatusString)
	}

	// OK
	return &reply.TargetLibQueryReplyData, nil
}
