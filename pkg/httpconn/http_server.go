package httpconn

import (
	"net"
	"time"
)

// HttpServer 基于Socket连接的HTTP服务端
type HttpServer struct {
	conn         net.Conn      // Socket连接通道
	readTimeout  time.Duration // 消息读取超时时间
	writeTimeout time.Duration // 消息发送超时时间

	priProtoHead *PrivateProtocolHead // 私有协议头配置
}

// NewHttpServer 创建基于Socket连接的HTTP服务端
func NewHttpServer(conn net.Conn) *HttpServer {
	return &HttpServer{
		conn:         conn,
		readTimeout:  time.Second * 30,
		writeTimeout: time.Second * 30,
	}
}

// SetReadTimeout 设置请求响应超时时间
func (s *HttpServer) SetReadTimeout(readTimeout, writeTimeout time.Duration) *HttpServer {
	s.readTimeout = readTimeout
	s.writeTimeout = writeTimeout
	return s
}

// SetPrivateProtocolHead 设置私有协议头
func (s *HttpServer) SetPrivateProtocolHead(opt PrivateProtocolHead) *HttpServer {
	// 克隆协议
	s.priProtoHead = opt.Clone()
	// OK
	return s
}

// Writer 获取响应对象
func (s *HttpServer) Writer() *HttpServerResponse {
	return NewHttpServerResponse(s)
}

// Reader 获取请求对象
func (s *HttpServer) Reader() *HttpServerRequest {
	return NewHttpServerRequest(s)
}
