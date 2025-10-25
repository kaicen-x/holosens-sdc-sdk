/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备连接托管
 */
package holosenssdcsdk

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/kaicen-x/holosens-sdc-sdk/api/application/device"
	"github.com/kaicen-x/holosens-sdc-sdk/api/application/metadata"
	"github.com/kaicen-x/holosens-sdc-sdk/api/application/snapshot"
	"github.com/kaicen-x/holosens-sdc-sdk/api/details/itgt"
	"github.com/kaicen-x/holosens-sdc-sdk/pkg/httpconn"
)

// Session 会话
type Session struct {
	connInstance    *httpconn.Connect // Socket连接实例
	deviceManager   *device.Manager   // 设备管理与维护管理器
	metadataManager *metadata.Manager // 智能元数据对接管理器
	snapshotManager *snapshot.Manager // 抓拍与图片下载管理器
	itgtManager     *itgt.Manager     // 智能分析管理器
}

// 新建会话
//
//	@param conn: 设备Socket连接通道
//	@return 会话
//	@return 错误信息
func newSession(conn net.Conn) *Session {
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
	return &Session{
		connInstance:    connInstance,
		deviceManager:   device.NewManager(connInstance),
		metadataManager: metadata.NewManager(connInstance),
		snapshotManager: snapshot.NewManager(connInstance),
		itgtManager:     itgt.NewManager(connInstance),
	}
}

// Close 关闭会话（同时会关闭Socket连接）
func (p *Session) Close() {
	p.connInstance.Close()
}

// GetHttp 获取基于Socket连接的HTTP会话
func (p *Session) GetHttp() *httpconn.Connect {
	return p.connInstance
}

// IsSetAuthorization 是否已设置认证信息
func (p *Session) IsSetAuthorization() bool {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	return client.IsSetAuthorization()
}

// BindAuthorizationChangeEvent 绑定认证信息修改事件
func (p *Session) BindAuthorizationChangeEvent(callback func(isClear bool)) {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	client.BindAuthorizationChangeEvent(callback)
}

// SetAuthorization 设置连接认证信息
func (p *Session) SetAuthorization(username, password string) {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	client.SetDigestAuth(username, password)
}

// DeviceManager 获取设备管理与维护管理器
func (p *Session) DeviceManager() *device.Manager {
	return p.deviceManager
}

// MetadataManager 获取智能元数据对接管理器
func (p *Session) MetadataManager() *metadata.Manager {
	return p.metadataManager
}

// SnapshotManager 获取抓拍与图片下载管理器
func (p *Session) SnapshotManager() *snapshot.Manager {
	return p.snapshotManager
}

// ItgtManager 获取智能分析管理器
func (p *Session) ItgtManager() *itgt.Manager {
	return p.itgtManager
}
