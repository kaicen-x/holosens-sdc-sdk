/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备主动注册服务端实现
 */
package holosenssdcsdk

import (
	"crypto/tls"
	"net"
	"time"

	devicemanage "github.com/bearki/holosens-sdc-sdk/api/device-manage"
	intelligentmetadata "github.com/bearki/holosens-sdc-sdk/api/intelligent-metadata"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// DeviceInstance 设备实例
type DeviceInstance struct {
	connInstance             *httpconn.Connect                     // Socket连接实例
	InitiativeRegisterParams devicemanage.InitiativeRegisterParams // 设备主动注册参数
	deviceManager            *devicemanage.Manager                 // 设备管理与维护管理器
	metadataManager          *intelligentmetadata.Manager          // 智能元数据对接管理器
}

// NewDeviceConnect 新建设备连接
//
//	@param conn: 设备Socket连接通道
//	@return 设备实例
//	@return 错误信息
func NewDeviceConnect(conn net.Conn) (*DeviceInstance, error) {
	// 处理连接保活
	switch tmpConn := conn.(type) {
	// 普通TCP连接
	case *net.TCPConn:
		tmpConn.SetKeepAlive(true)
		tmpConn.SetKeepAlivePeriod(time.Minute)
	// TLS连接
	case *tls.Conn:
		tmpConn.NetConn().(*net.TCPConn).SetKeepAlive(true)
		tmpConn.NetConn().(*net.TCPConn).SetKeepAlivePeriod(time.Minute)
	}

	// 构建连接实例
	connInstance := httpconn.NewConnect(true, conn)
	// 创建设备管理与维护管理器
	deviceManager := devicemanage.NewManager(connInstance)
	// 创建智能元数据对接管理器
	metadataManager := intelligentmetadata.NewManager(connInstance)

	// 接收设备主动注册信息
	params, err := deviceManager.InitiativeRegister()
	if err != nil {
		return nil, err
	}

	// 返回设备实例
	return &DeviceInstance{
		connInstance:             connInstance,
		InitiativeRegisterParams: *params,
		deviceManager:            deviceManager,
		metadataManager:          metadataManager,
	}, nil
}

// SetTimeout 设置连接读取超时时间
func (p *DeviceInstance) SetReadTimeout(timeout time.Duration) {
	p.connInstance.SetReadTimeout(timeout)
}

// SetTimeout 设置连接写入超时时间
func (p *DeviceInstance) SetWriteTimeout(timeout time.Duration) {
	p.connInstance.SetWriteTimeout(timeout)
}

// SetAuthorization 设置连接认证信息
func (p *DeviceInstance) SetAuthorization(username, password string) {
	p.connInstance.SetAuthorization(username, password)
}

// DeviceManager 获取设备管理与维护管理器
func (p *DeviceInstance) DeviceManager() *devicemanage.Manager {
	return p.deviceManager
}

// IntelligentMetadataManager 获取智能元数据对接管理器
func (p *DeviceInstance) IntelligentMetadataManager() *intelligentmetadata.Manager {
	return p.metadataManager
}
