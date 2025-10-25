/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标记录添加
 */
package recognize

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
)

// TargetRecordCreateParams 目标记录添加参数
type TargetRecordCreateParams struct {
	// 目标库信息（必填）
	Lib TargetRecordParamsWithLib `json:"facelib"`
	// 目标记录创建信息（必填）
	Record TargetRecordCreateInfo `json:"record"`
}

// TargetRecordCreateReply 目标记录添加响应
type TargetRecordCreateReply struct {
	// 目标库简要信息
	Lib TargetLibBriefInfo `json:"facelib"`
	// 目标记录信息
	Record TargetRecordInfo `json:"record"`
}

// 填充目标记录创建表单
func fillTargetRecordCreateFormData(formData *multipart.Writer, params *TargetRecordCreateParams, img []byte) error {
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

// TargetRecordCreate 目标记录添加
//
//	@param	params: 目标记录添加参数
//	@param	img: 目标记录图片
//	@return	目标记录添加响应
//	@return	错误信息
func (p *Manager) TargetRecordCreate(params TargetRecordCreateParams, img []byte) (*TargetRecordCreateReply, error) {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 构建表单数据
	formBuf := new(bytes.Buffer)
	formData := multipart.NewWriter(formBuf)
	// 填充表单
	err := fillTargetRecordCreateFormData(formData, &params, img)
	if err != nil {
		formData.Close()
		return nil, err
	}
	// 闭合表单
	err = formData.Close()
	if err != nil {
		return nil, err
	}

	// 发送请求
	var reply TargetRecordCreateReply
	_, err = client.Post("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/FaceRecord").
		SetQuery("TaskType", "1").
		SetBody(io.NopCloser(formBuf), int64(formBuf.Len())).
		SetContentType(formData.FormDataContentType()).
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
