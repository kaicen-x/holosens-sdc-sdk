package httpconn

import (
	"encoding/json"
	"io"
	"net/http"
)

// HttpServerRequest HTTP自定义服务端请求对象
type HttpServerRequest struct {
	ser *HttpServer   // 服务端
	req *http.Request // 原始请求对象
}

// NewHttpServerRequest 创建HTTP自定义服务端请求对象
func NewHttpServerRequest(server *HttpServer) *HttpServerRequest {
	// 构建HTTP自定义服务端请求对象
	return &HttpServerRequest{
		ser: server,
	}
}

// RawRequest 获取原始请求对象
func (r *HttpServerRequest) RawRequest() *http.Request {
	// 是否存在缓存
	if r.req != nil {
		return r.req
	}
	// 读取HTTP请求
	req, err := readHttpRequest(r.ser)
	if err != nil {
		return nil
	}
	// 缓存原始请求对象
	r.req = req
	// OK
	return req
}

// BindJSON 绑定JSON数据（请不要重复Close Body）
func (r *HttpServerRequest) BindJSON(obj any) error {
	// 读取HTTP请求
	req, err := readHttpRequest(r.ser)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	// 缓存原始请求对象
	r.req = req
	// 读取请求体
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	// 反序列化JSON
	return json.Unmarshal(body, obj)
}
