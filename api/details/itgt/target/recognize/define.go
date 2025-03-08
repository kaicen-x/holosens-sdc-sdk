/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标识别模块
 */
package recognize

import (
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 目标识别管理器
type Manager struct {
	connInstance *httpconn.Connect // Socket连接实例
}

// NewManager 创建目标识别管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance: connInstance,
	}
}
