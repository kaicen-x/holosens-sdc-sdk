/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备主动注册（服务端）Socket连接托管
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

// DeviceConnectWithServer 设备Socket连接服务端托管
type DeviceConnectWithServer struct {
	*DeviceConnect                                                 // 设备Socket连接托管基础
	InitiativeRegisterParams devicemanage.InitiativeRegisterParams // 设备主动注册参数
}

// NewDeviceConnectWithServer 新建设备连接实例（来自服务端）
//
//	@param conn: 设备Socket连接通道
//	@return 设备实例
//	@return 错误信息
func NewDeviceConnectWithServer(conn net.Conn) (*DeviceConnectWithServer, error) {
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
	// 获取连接实例上的客户端
	client := connInstance.LockHttpClient()
	// 设置客户端私有响应协议
	client.SetPrivateProtocolHead(httpconn.PrivateProtocolHead{
		ResponseHead: make([]byte, 8),
		Strict:       false,
	})
	// 释放连接实例占用
	connInstance.Unlock()

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
	return &DeviceConnectWithServer{
		DeviceConnect: &DeviceConnect{
			connInstance:    connInstance,
			deviceManager:   deviceManager,
			metadataManager: metadataManager,
		},
		InitiativeRegisterParams: *params,
	}, nil
}
