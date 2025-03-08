/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备连接托管
 */
package holosenssdcsdk

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/bearki/holosens-sdc-sdk/api/application/device"
	"github.com/bearki/holosens-sdc-sdk/api/application/metadata"
	"github.com/bearki/holosens-sdc-sdk/api/application/snapshot"
	"github.com/bearki/holosens-sdc-sdk/api/details/itgt"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// ConnectTrusteeship Socket连接托管器
type ConnectTrusteeship struct {
	connInstance    *httpconn.Connect // Socket连接实例
	deviceManager   *device.Manager   // 设备管理与维护管理器
	metadataManager *metadata.Manager // 智能元数据对接管理器
	snapshotManager *snapshot.Manager // 抓拍与图片下载管理器
	itgtManager     *itgt.Manager     // 智能分析管理器
}

// NewWithClient 新建客户端Socket连接托管器
//
//	@param conn: 设备Socket连接通道
//	@return 设备实例
//	@return 错误信息
func newConnectTrusteeship(conn net.Conn) *ConnectTrusteeship {
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
	// OK
	return &ConnectTrusteeship{
		connInstance:    connInstance,
		deviceManager:   device.NewManager(connInstance),
		metadataManager: metadata.NewManager(connInstance),
		snapshotManager: snapshot.NewManager(connInstance),
		itgtManager:     itgt.NewManager(connInstance),
	}
}

// GetHttp 获取基于Socket连接的HTTP托管器
func (p *ConnectTrusteeship) GetHttp() *httpconn.Connect {
	return p.connInstance
}

// SetAuthorization 设置连接认证信息
func (p *ConnectTrusteeship) SetAuthorization(username, password string) {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	client.SetDigestAuth(username, password)
}

// DeviceManager 获取设备管理与维护管理器
func (p *ConnectTrusteeship) DeviceManager() *device.Manager {
	return p.deviceManager
}

// MetadataManager 获取智能元数据对接管理器
func (p *ConnectTrusteeship) MetadataManager() *metadata.Manager {
	return p.metadataManager
}

// SnapshotManager 获取抓拍与图片下载管理器
func (p *ConnectTrusteeship) SnapshotManager() *snapshot.Manager {
	return p.snapshotManager
}

// ItgtManager 获取智能分析管理器
func (p *ConnectTrusteeship) ItgtManager() *itgt.Manager {
	return p.itgtManager
}
