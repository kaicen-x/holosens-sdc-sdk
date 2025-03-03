package httpconn

import "time"

// HttpServer 基于Socket连接的HTTP服务端
type HttpServer struct {
	_connInfo    *ConnectInfo  // Socket连接信息
	readTimeout  time.Duration // 请求响应超时时间
	writeTimeout time.Duration // 请求发送超时时间
}

// NewHttpServer 创建基于Socket连接的HTTP服务端
func NewHttpServer(cii *ConnectInfo) *HttpServer {
	return &HttpServer{
		_connInfo:    cii,
		readTimeout:  0,
		writeTimeout: 0,
	}
}

// SetReadTimeout 设置请求响应超时时间
func (s *HttpServer) SetReadTimeout(readTimeout, writeTimeout time.Duration) *HttpServer {
	s.readTimeout = readTimeout
	s.writeTimeout = writeTimeout
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
