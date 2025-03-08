/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能分析模块
 */
package itgt

import (
	"github.com/bearki/holosens-sdc-sdk/api/details/itgt/target"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 智能分析管理器
type Manager struct {
	connInstance  *httpconn.Connect // Socket连接实例
	targetManager *target.Manager   // 目标相关管理器
}

// NewManager 创建智能分析管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance:  connInstance,
		targetManager: target.NewManager(connInstance),
	}
}

// TargetManager 获取目标相关管理器
func (m *Manager) TargetManager() *target.Manager {
	return m.targetManager
}
