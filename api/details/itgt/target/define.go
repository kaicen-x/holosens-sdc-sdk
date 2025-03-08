/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标相关模块
 */
package target

import (
	"github.com/bearki/holosens-sdc-sdk/api/details/itgt/target/recognize"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// Manager 目标相关管理器
type Manager struct {
	connInstance     *httpconn.Connect  // Socket连接实例
	recognizeManager *recognize.Manager // 目标识别管理器
}

// NewManager 创建目标相关管理器
func NewManager(connInstance *httpconn.Connect) *Manager {
	return &Manager{
		connInstance:     connInstance,
		recognizeManager: recognize.NewManager(connInstance),
	}
}

// RecognizeManager 获取目标识别管理器
func (p *Manager) RecognizeManager() *recognize.Manager {
	return p.recognizeManager
}
