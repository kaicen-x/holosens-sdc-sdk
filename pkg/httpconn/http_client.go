package httpconn

import (
	"net/http"
	"time"
)

// HttpClient 基于Socket连接的HTTP客户端
type HttpClient struct {
	_connInfo    *ConnectInfo  // Socket连接信息
	readTimeout  time.Duration // 请求响应超时时间
	writeTimeout time.Duration // 请求发送超时时间
}

// NewHttpClient 创建基于Socket连接的HTTP客户端
func NewHttpClient(cii *ConnectInfo) *HttpClient {
	return &HttpClient{
		_connInfo:    cii,
		readTimeout:  0,
		writeTimeout: 0,
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
