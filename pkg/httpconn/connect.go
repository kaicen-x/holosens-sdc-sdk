package httpconn

import (
	"net"
	"sync"
	"time"
)

// ConnectInfo Socket连接实例信息
type ConnectInfo struct {
	connType         bool          // 连接类型（true: 主动注册，false: 被动注册）
	conn             net.Conn      // Socket连接通道
	connReadTimeout  time.Duration // Socket连接通道读取超时时间
	connWriteTimeout time.Duration // Socket连接通道写入超时时间
	username         string        // 认证用户名
	password         string        // 认证密码
}

// HttpClient 构建HTTP客户端
func (cii *ConnectInfo) HttpClient() *HttpClient {
	return NewHttpClient(cii)
}

// HttpServer 获取HTTP服务端
func (cii *ConnectInfo) HttpServer() *HttpServer {
	return NewHttpServer(cii)
}

// Connect 连接实例
type Connect struct {
	connMtx  sync.Mutex   // 连接锁（每个设备全局仅有一个）
	connInfo *ConnectInfo // 连接实例信息
}

// NewConnect 创建连接实例
//
//	@param conn: Socket连接通道
//	@param connType: 连接类型（true: 主动注册，false: 被动注册）
//	@return 连接实例
func NewConnect(connType bool, conn net.Conn) *Connect {
	return &Connect{
		connInfo: &ConnectInfo{
			connType:         connType,
			conn:             conn,
			connReadTimeout:  time.Second * 30,
			connWriteTimeout: time.Second * 30,
		},
	}
}

// Lock 获取连接锁
//
//	@return 连接实例信息
func (ci *Connect) Lock() *ConnectInfo {
	ci.connMtx.Lock()
	return ci.connInfo
}

// Unlock 释放连接锁
func (ci *Connect) Unlock() {
	ci.connMtx.Unlock()
}

// SetTimeout 设置连接读取超时时间
func (ci *Connect) SetReadTimeout(timeout time.Duration) {
	connInfo := ci.Lock()
	defer ci.Unlock()
	connInfo.connReadTimeout = timeout
}

// SetTimeout 设置连接写入超时时间
func (ci *Connect) SetWriteTimeout(timeout time.Duration) {
	connInfo := ci.Lock()
	defer ci.Unlock()
	connInfo.connWriteTimeout = timeout
}

// SetAuthorization 设置连接认证信息
func (ci *Connect) SetAuthorization(username, password string) {
	connInfo := ci.Lock()
	defer ci.Unlock()
	connInfo.username = username
	connInfo.password = password
}
