/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标记录批量删除
 */
package recognize

import (
	"errors"

	"github.com/kaicen-x/holosens-sdc-sdk/api/common"
)

// TargetRecordBatchDeleteParams 目标记录批量删除参数
//
//	1、全部删除，gender=-1，cardType=-1，isStore=-1，其他字段，数字为0，字符串的为空
//	2、条件删除，未填的字段：除以上三个字段外，其他类型数字为0，字符串的为空
type TargetRecordBatchDeleteParams struct {
	// 目标库信息（必填）
	Lib TargetRecordParamsWithLib `json:"facelib"`
	// 要删除的记录总数（必填）
	Num int `json:"faceNum"`
	// 要删除的目标记录ID列表（必填）
	IDs []int `json:"faceId"`
}

// TargetRecordBatchDeleteReply 目标记录批量删除响应
type TargetRecordBatchDeleteReply = common.Response[common.ResponseStatus]

// TargetRecordBatchDelete 目标记录批量删除
//
//	@param	params: 目标记录批量删除参数
//	@return	错误信息
func (p *Manager) TargetRecordBatchDelete(params TargetRecordBatchDeleteParams) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetRecordBatchDeleteReply
	_, err := client.Delete("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/FaceRecord").
		SetQuery("TaskType", "1").
		SetJSON(&params).
		SetContentType("application/json").
		DecodeJSON(&reply)
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
