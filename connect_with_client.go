/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备被动注册（客户端）Socket连接托管
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

// DeviceConnectWithClient 设备Socket连接服务端托管
type DeviceConnectWithClient struct {
	*DeviceConnect // 设备Socket连接托管基础
}

// NewDeviceConnectWithClient 新建设备连接实例（来自客户端）
//
//	@param conn: 设备Socket连接通道
//	@return 设备实例
//	@return 错误信息
func NewDeviceConnectWithClient(conn net.Conn) *DeviceConnectWithClient {
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
	connInstance := httpconn.NewConnect(conn)
	// 创建设备管理与维护管理器
	deviceManager := devicemanage.NewManager(connInstance)
	// 创建智能元数据对接管理器
	metadataManager := intelligentmetadata.NewManager(connInstance)

	// 返回设备实例
	return &DeviceConnectWithClient{
		DeviceConnect: &DeviceConnect{
			connInstance:    connInstance,
			deviceManager:   deviceManager,
			metadataManager: metadataManager,
		},
	}
}
