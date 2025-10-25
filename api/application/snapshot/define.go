/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口抓拍与图片下载模块
 */
package snapshot

import (
	"github.com/kaicen-x/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 抓拍与图片下载管理器
type Manager struct {
	connInstance *httpconn.Connect // Socket连接实例
}

// NewManager 创建抓拍与图片下载管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance: connInstance,
	}
}
