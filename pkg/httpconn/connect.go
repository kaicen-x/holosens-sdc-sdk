package httpconn

import (
	"net"
	"sync"
)

// Connect 连接实例
type Connect struct {
	mtx    sync.Mutex  // 连接锁（HTTP客户端和HTTP服务端同时只能使用一个）
	conn   net.Conn    // Socket连接通道
	client *HttpClient // Socket连接通道上的HTTP客户端
	server *HttpServer // Socket连接通道上的HTTP服务端
}

// NewConnect 创建连接实例
//
//	@param conn: Socket连接通道
//	@return 连接实例
func NewConnect(conn net.Conn) *Connect {
	return &Connect{
		conn:   conn,
		client: NewHttpClient(conn),
		server: NewHttpServer(conn),
	}
}

// HttpClient 获取HTTP客户端
func (ci *Connect) LockHttpClient() *HttpClient {
	ci.mtx.Lock()
	return ci.client
}

// HttpServer 获取HTTP服务端
func (ci *Connect) LockHttpServer() *HttpServer {
	ci.mtx.Lock()
	return ci.server
}

// Unlock 释放连接锁
//
//	调用完以下接口后请及时调用该接口释放Socket占用，保证连接可及时供他人获取
//	LockHttpClient()
//	LockHttpServer()
func (ci *Connect) Unlock() {
	ci.mtx.Unlock()
}
