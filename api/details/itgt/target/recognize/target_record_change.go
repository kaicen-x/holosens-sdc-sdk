/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标记录修改
 */
package recognize

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"

	"github.com/bearki/holosens-sdc-sdk/api/common"
)

// TargetRecordChangeParams 目标记录修改参数
type TargetRecordChangeParams struct {
	// 目标库信息（必填）
	Lib TargetRecordParamsWithLib `json:"facelib"`
	// 目标记录信息（必填）
	Record TargetRecordChangeInfo `json:"record"`
}

// TargetRecordChangeReply 目标记录修改响应
type TargetRecordChangeReply = common.Response[common.ResponseStatus]

// 填充目标记录修改表单
func fillTargetRecordChangeFormData(formData *multipart.Writer, params *TargetRecordChangeParams, img []byte) error {
	// 创建FILE字段
	fileWriter, err := formData.CreateFormFile("file", params.Record.PicName)
	if err != nil {
		return err
	}
	// 写入文件内容
	fileCount, err := fileWriter.Write(img)
	if err != nil {
		return err
	}
	if fileCount != len(img) {
		return errors.New("file write failed, write count not match")
	}

	// 序列化JSON数据
	jsonData, err := json.Marshal(&params)
	if err != nil {
		return err
	}
	// 创建JSON字段
	jsonWriter, err := formData.CreateFormField("json")
	if err != nil {
		return err
	}
	// 写入JSON数据
	jsonCount, err := jsonWriter.Write(jsonData)
	if err != nil {
		return err
	}
	if jsonCount != len(jsonData) {
		return errors.New("json data write error, write count not match")
	}

	// OK
	return nil
}

// TargetRecordChange 目标记录修改
//
//	@param	params: 目标记录修改参数
//	@param	img: 目标记录图片
//	@return	错误信息
func (p *Manager) TargetRecordChange(params TargetRecordChangeParams, img []byte) error {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 构建表单数据
	formBuf := new(bytes.Buffer)
	formData := multipart.NewWriter(formBuf)
	// 填充表单
	err := fillTargetRecordChangeFormData(formData, &params, img)
	if err != nil {
		formData.Close()
		return err
	}
	// 闭合表单
	err = formData.Close()
	if err != nil {
		return err
	}

	// 发送请求
	var reply TargetRecordChangeReply
	_, err = client.Put("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/FaceRecord").
		SetQuery("TaskType", "2").
		SetBody(io.NopCloser(formBuf), int64(formBuf.Len())).
		SetContentType(formData.FormDataContentType()).
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
