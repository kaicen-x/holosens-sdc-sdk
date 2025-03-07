/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备管理与维护模块
 */
package device

import (
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 设备管理与维护管理器
type Manager struct {
	connInstance *httpconn.Connect // Socket连接实例
}

// NewManager 创建设备管理与维护管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance: connInstance,
	}
}
