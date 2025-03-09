/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标库新建
 */
package recognize

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// TargetLibCreateParams 目标库新建参数
type TargetLibCreateParams struct {
	// 目标库基本信息（必填）
	FaceLib TargetLibBaseInfo `json:"facelib"`
}

// TargetLibCreateReply 目标库新建响应
type TargetLibCreateReply = common.Response[common.ResponseStatus]

// TargetLibCreate 目标库新建
//
//	@param	params: 目标库新建参数
//	@return 异常信息
func (p *Manager) TargetLibCreate(params TargetLibCreateParams) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetLibCreateReply
	_, err := client.Post("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/Libs").
		SetJSON(&params).
		SetContentType("application/json").
		DecodeJSON(&reply)
	if err != nil {
		return err
	}

	// 检查是否创建成功
	if reply.ResponseStatus.StatusCode != 0 {
		return errors.New(reply.ResponseStatus.StatusString)
	}

	// OK
	return nil
}
