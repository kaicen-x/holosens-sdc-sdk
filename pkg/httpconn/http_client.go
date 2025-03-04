package httpconn

import (
	"net"
	"net/http"
	"time"
)

// HttpClientAuthType HTTP客户端认证类型
type HttpClientAuthType int

// HTTP客户端认证类型美剧
const (
	HttpClientAuthTypeNone   HttpClientAuthType = iota // 无认证
	HttpClientAuthTypeBasic                            // 认证类型：基本认证
	HttpClientAuthTypeDigest                           // 认证类型：摘要认证
	HttpClientAuthTypeOAuth                            // 认证类型：OAuth认证；TODO: 未实现
)

// HttpClientAuth HTTP客户端认证信息
type HttpClientAuth struct {
	Type     HttpClientAuthType // 认证类型
	Username string             // Basic/Digest认证用户名
	Password string             // Basic/Digest认证密码
}

// HttpClient 基于Socket连接的HTTP客户端
type HttpClient struct {
	conn         net.Conn      // Socket连接通道
	readTimeout  time.Duration // 消息读取超时时间
	writeTimeout time.Duration // 消息发送超时时间

	auth         *HttpClientAuth      // 认证信息
	priProtoHead *PrivateProtocolHead // 私有协议头配置
}

// NewHttpClient 创建基于Socket连接的HTTP客户端
func NewHttpClient(conn net.Conn) *HttpClient {
	return &HttpClient{
		conn:         conn,
		readTimeout:  time.Second * 30,
		writeTimeout: time.Second * 30,
		auth:         nil,
		priProtoHead: nil,
	}
}

// SetTimeout 设置请求响应超时时间
func (c *HttpClient) SetTimeout(readTimeout, writeTimeout time.Duration) *HttpClient {
	// 设置请求响应超时时间
	c.readTimeout = readTimeout
	// 设置请求发送超时时间
	c.writeTimeout = writeTimeout
	// OK
	return c
}

// SetDigestAuth 设置Digest认证信息
func (c *HttpClient) SetDigestAuth(username, password string) *HttpClient {
	c.auth = &HttpClientAuth{
		Type:     HttpClientAuthTypeDigest,
		Username: username,
		Password: password,
	}
	// OK
	return c
}

// SetBasicAuth 设置Basic认证信息
func (c *HttpClient) SetBasicAuth(username string, password ...string) *HttpClient {
	c.auth = &HttpClientAuth{
		Type:     HttpClientAuthTypeBasic,
		Username: username,
		Password: "",
	}
	if len(password) > 0 {
		c.auth.Password = password[0]
	}
	// OK
	return c
}

// SetPrivateProtocolHead 设置私有协议头
func (c *HttpClient) SetPrivateProtocolHead(opt PrivateProtocolHead) *HttpClient {
	// 克隆协议
	c.priProtoHead = opt.Clone()
	// OK
	return c
}

// Get 封装Get请求
func (c *HttpClient) Get(url string) *HttpClientRequest {
	return NewHttpClientRequest(c, http.MethodGet, url)
}

// Post 封装Post请求
func (c *HttpClient) Post(url string) *HttpClientRequest {
	return NewHttpClientRequest(c, http.MethodPost, url)
}

// Put 封装Put请求
func (c *HttpClient) Put(url string) *HttpClientRequest {
	return NewHttpClientRequest(c, http.MethodPut, url)
}

// Delete 封装Delete请求
func (c *HttpClient) Delete(url string) *HttpClientRequest {
	return NewHttpClientRequest(c, http.MethodDelete, url)
}
