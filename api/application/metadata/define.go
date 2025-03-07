/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据对接模块
 */
package metadata

import (
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 智能元数据对接管理器
type Manager struct {
	connInstance *httpconn.Connect // Socket连接实例
}

// NewManager 创建智能元数据对接管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance: connInstance,
	}
}
