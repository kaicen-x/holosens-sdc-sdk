/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口手动抓拍
 */
package snapshot

import (
	"errors"
	"io"
	"strings"
)

// SnapActionParams 手动抓拍参数
type SnapActionParams struct {
	UUID      string `json:"UUID"`      // 通道UUID
	ChannelID int    `json:"channelID"` // 通道ID（选填）
}

// SnapActionReply 手动抓拍响应
type SnapActionReply struct {
	Data        []byte // 抓拍图片数据
	ContentType string // 抓拍图片格式
	FileName    string // 抓拍图片文件名
}

// SnapAction 手动抓拍
//
//	@param params: 手动抓拍参数
//	@return: 手动抓拍响应
//	@return: 错误信息
func (p *Manager) SnapAction(params SnapActionParams) (*SnapActionReply, error) {
	// 获取Socket连接的HTTP客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	form, _, err := client.Post("/SDCAPI/V1.0/Storage/Snapshot/SnapAction").
		SetJSON(&params).
		DecodeFormData(1024 * 1024 * 5)
	if err != nil {
		return nil, err
	}
	defer form.RemoveAll()

	// 读取响应
	// 遍历所有表单键
	for _, fs := range form.File {
		// 遍历所有表单值
		for _, fh := range fs {
			// 跳过非图片文件
			if !strings.Contains(fh.Header.Get("Content-Type"), "jpeg") {
				continue
			}

			// 打开图片文件
			file, err := fh.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			// 读取图片数据
			data, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			// OK
			return &SnapActionReply{
				Data:        data,
				ContentType: fh.Header.Get("Content-Type"),
				FileName:    fh.Filename,
			}, nil
		}
	}

	// 错误
	return nil, errors.New("no file")
}
