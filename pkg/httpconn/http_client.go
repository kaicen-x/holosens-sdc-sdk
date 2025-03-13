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

	auth            *HttpClientAuth      // 认证信息
	authChangeEvent func(isClear bool)   // 认证信息修改事件
	priProtoHead    *PrivateProtocolHead // 私有协议头配置
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

// IsSetAuthorization 是否已设置认证信息
func (p *HttpClient) IsSetAuthorization() bool {
	return p.auth != nil && p.auth.Type != HttpClientAuthTypeNone
}

// BindAuthorizationChangeEvent 绑定认证信息修改事件
func (p *HttpClient) BindAuthorizationChangeEvent(callback func(isClear bool)) {
	p.authChangeEvent = callback
}

// 设置认证信息
func (c *HttpClient) setAuthorization(authType HttpClientAuthType, user, pass string) *HttpClient {
	// 是否需要回调
	if c.authChangeEvent != nil && (c.auth == nil ||
		c.auth.Type != authType ||
		c.auth.Username != user ||
		c.auth.Password != pass) {
		// 触发回调
		c.authChangeEvent(false)
		// 赋值新授权信息
		c.auth = &HttpClientAuth{
			Type:     authType,
			Username: user,
			Password: pass,
		}
	}
	// OK
	return c
}

// SetDigestAuth 设置Digest认证信息
func (c *HttpClient) SetDigestAuth(username, password string) *HttpClient {
	return c.setAuthorization(HttpClientAuthTypeDigest, username, password)
}

// SetBasicAuth 设置Basic认证信息
func (c *HttpClient) SetBasicAuth(username string, password ...string) *HttpClient {
	newPassword := ""
	if len(password) > 0 {
		newPassword = password[0]
	}
	return c.setAuthorization(HttpClientAuthTypeDigest, username, newPassword)
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
