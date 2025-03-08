/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标识别库修改
 */
package recognize

import (
	"errors"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// TargetLibChangeData 目标库修改数据
type TargetLibChangeData struct {
	// 目标库ID（必填）
	ID int `json:"faceLibId"`
	// 目标库设防状态（选填，默认：0）
	//	取值范围：0-未设防, 1-设防
	OnControl int `json:"control"`
	// 目标库基础信息
	TargetLibBaseInfo
}

// TargetLibChangeParams 目标库修改参数
type TargetLibChangeParams struct {
	// 目标库修改参数数据（必填）
	FaceLib TargetLibChangeData `json:"facelib"`
}

// TargetLibChangeReply 目标库修改响应
type TargetLibChangeReply = common.Response[common.ResponseStatus]

// TargetLibChange 目标库修改
//
//	@return 目标库修改参数
//	@return 异常信息
func (p *Manager) TargetLibChange(params TargetLibChangeParams) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetLibChangeReply
	_, err := client.Put("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/Libs").
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
