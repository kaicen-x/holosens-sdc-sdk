/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备连接托管
 */
package holosenssdcsdk

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/kaicen-x/holosens-sdc-sdk/api/application/device"
	"github.com/kaicen-x/holosens-sdc-sdk/api/application/metadata"
	"github.com/kaicen-x/holosens-sdc-sdk/api/application/snapshot"
	"github.com/kaicen-x/holosens-sdc-sdk/api/details/itgt"
	"github.com/kaicen-x/holosens-sdc-sdk/pkg/httpconn"
)

// Session 会话
type Session struct {
	httpConn        *httpconn.Connect // HTTP连接通道
	deviceManager   *device.Manager   // 设备管理与维护管理器
	metadataManager *metadata.Manager // 智能元数据对接管理器
	snapshotManager *snapshot.Manager // 抓拍与图片下载管理器
	itgtManager     *itgt.Manager     // 智能分析管理器
}

// 新建会话
//
//	@param conn: Socket连接通道
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

	// 构建HTTP连接通道实例
	httpConn := httpconn.NewConnect(conn)
	// OK
	return &Session{
		httpConn:        httpConn,
		deviceManager:   device.NewManager(httpConn),
		metadataManager: metadata.NewManager(httpConn),
		snapshotManager: snapshot.NewManager(httpConn),
		itgtManager:     itgt.NewManager(httpConn),
	}
}

// Close 关闭会话（同时会关闭Socket连接）
func (p *Session) Close() {
	p.httpConn.Close()
}

// GetHttp 获取HTTP连接通道
func (p *Session) GetHttp() *httpconn.Connect {
	return p.httpConn
}

// IsSetAuthorization 是否已设置认证信息
//
// 可用于防止重复设置认证信息
func (p *Session) IsSetAuthorization() bool {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	return client.IsSetAuthorization()
}

// BindAuthorizationChangeEvent 绑定认证信息修改事件
//
// 该方法由SessionCache内部使用，以便于认证信息切换后重新启动新的心跳检测，通常外部无需使用或关心该方法。
// 除非您要使用自定义的会话缓存器，那么该方法对您来说可能有用。
func (p *Session) BindAuthorizationChangeEvent(callback func(isClear bool)) {
	client := p.GetHttp().LockHttpClient()
	defer p.GetHttp().Unlock()
	client.BindAuthorizationChangeEvent(callback)
}

// SetAuthorization 设置连接认证信息
//
// 一般情况下北向接口需要设置认证信息才能与摄像头正常连接
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

// SessionWithServer 服务端会话
type SessionWithServer struct {
	*Session                                                 // 会话
	InitiativeRegisterParams device.InitiativeRegisterParams // 设备主动注册参数
}

// 设置私有协议头
//
// 在主动注册场景中，由于历史原因，摄像机返回的响应消息最开始有8个字节的私有协
// 议头，解析响应时需对此进行处理，否则可能导致解析出错。
func setPrivateProtocolHead(session *Session) {
	// 获取连接实例上的客户端
	client := session.GetHttp().LockHttpClient()
	defer session.GetHttp().Unlock()
	// 设置HTTP客户端私有协议头
	client.SetPrivateProtocolHead(httpconn.PrivateProtocolHead{
		RequestHead:  nil,
		ResponseHead: make([]byte, 8),
		Strict:       false,
	})
}

// NewWithTcpServer 托管服务端会话（基于TCP服务器）
//
//	@param conn: 设备TCP连接通道
//	@return 服务端会话
//	@return 错误信息
func NewWithTcpServer(conn net.Conn) (*SessionWithServer, error) {
	// 创建会话
	session := newSession(conn)
	// 设置私有协议头
	setPrivateProtocolHead(session)

	// 接收设备主动注册信息
	params, err := session.DeviceManager().InitiativeRegister()
	if err != nil {
		return nil, err
	}

	// 创建并返回服务端会话
	return &SessionWithServer{
		Session:                  session,
		InitiativeRegisterParams: *params,
	}, nil
}

// NewWithHttpServer 托管服务端会话（基于HTTP服务器）
//
//	@param w: 设备HTTP请求响应对象写入器
//	@param r: 设备HTTP请求对象
//	@return 服务端会话
//	@return 错误信息
func NewWithHttpServer(w http.ResponseWriter, r *http.Request) (*SessionWithServer, error) {
	// 处理请求
	params, err := device.InitiativeRegister(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// 接管HTTP的TCP连接
	resController := http.NewResponseController(w)
	conn, _, err := resController.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// 创建会话
	session := newSession(conn)
	// 设置私有协议头
	setPrivateProtocolHead(session)

	// 创建并返回服务端会话
	return &SessionWithServer{
		Session:                  session,
		InitiativeRegisterParams: *params,
	}, nil
}

// SessionWithClient 客户端会话
type SessionWithClient struct {
	*Session // 会话
}

// NewWithTcpClient 托管客户端会话
//
//	@param conn: 设备TCP连接通道
//	@return 客户端会话
//	@return 错误信息
func NewWithTcpClient(conn net.Conn) *SessionWithClient {
	// 返回客户端会话
	return &SessionWithClient{
		// 创建会话
		Session: newSession(conn),
	}
}
