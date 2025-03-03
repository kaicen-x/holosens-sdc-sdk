package httpconn

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// HttpServerResponse HTTP自定义服务端响应对象
type HttpServerResponse struct {
	_server   *HttpServer    // 服务端
	_cacheErr error          // 缓存链式调用中间产生的异常
	res       *http.Response // HTTP原生响应对象
}

// NewHttpServerResponse 创建HTTP自定义服务端响应对象
func NewHttpServerResponse(server *HttpServer) *HttpServerResponse {
	// 构建HTTP自定义服务端响应对象
	return &HttpServerResponse{
		_server:   server,
		_cacheErr: nil,
		res: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Connection":   []string{"keep-alive"},
				"Content-Type": []string{"application/json; charset=UTF-8"},
			},
			Body:          nil,
			ContentLength: 0,
		},
	}
}

// SetHeader 设置响应头
func (w *HttpServerResponse) SetHeader(key, val string) *HttpServerResponse {
	// 检查
	if w.res == nil {
		return w
	}
	// 设置响应头
	if w.res.Header == nil {
		w.res.Header = make(http.Header)
	}
	w.res.Header.Set(key, val)
	// OK
	return w
}

// AddHeader 添加响应头
func (w *HttpServerResponse) AddHeader(key, val string) *HttpServerResponse {
	// 检查
	if w.res == nil {
		return w
	}
	// 添加响应头
	if w.res.Header == nil {
		w.res.Header = make(http.Header)
	}
	w.res.Header.Add(key, val)
	// OK
	return w
}

// SetContentType 设置Content-Type响应头信息
func (w *HttpServerResponse) SetContentType(contentType string) *HttpServerResponse {
	// 检查
	if w.res == nil {
		return w
	}
	// 设置响应头信息
	if w.res.Header == nil {
		w.res.Header = make(http.Header)
	}
	w.res.Header.Set("Content-Type", contentType)
	// OK
	return w
}

// SetStatusCode 设置响应状态码
func (w *HttpServerResponse) SetStatusCode(statusCode int) *HttpServerResponse {
	// 检查
	if w.res == nil {
		return w
	}
	// 设置响应头信息
	w.res.StatusCode = statusCode
	// OK
	return w
}

// SetBody 设置响应体
func (w *HttpServerResponse) setBody(body io.ReadCloser, bodySize int64) *HttpServerResponse {
	// 检查
	if w.res == nil {
		body.Close() // 提前结束需要关闭响应体
		return w
	}
	// 设置响应体
	w.res.Body = body
	w.res.ContentLength = bodySize
	// OK
	return w
}

// SetJsonBody 设置JSON响应体
func (w *HttpServerResponse) setJSON(data any) *HttpServerResponse {
	// 检查
	if w.res == nil {
		return w
	}
	// json编码数据
	dataBytes, err := json.Marshal(data)
	if err != nil {
		w._cacheErr = errors.Join(w._cacheErr, err)
		return w
	}
	// 设置响应体
	w.setBody(io.NopCloser(bytes.NewReader(dataBytes)), int64(len(dataBytes)))
	// 设置响应头信息
	w.SetContentType("application/json; charset=UTF-8")
	// OK
	return w
}

// 执行响应
func (w *HttpServerResponse) send() error {
	// 提取原生响应对象
	res := w.res
	if res == nil {
		return errors.Join(w._cacheErr, errors.New("response is nil"))
	}
	if w._cacheErr != nil {
		if res.Body != nil {
			res.Body.Close() // 提前结束需要手动释放Body
		}
		return w._cacheErr
	}
	// 发送响应
	return writeHttpResponse(w._server._connInfo, res)
}

// Send 执行响应（请主动关闭Response.Body）
func (w *HttpServerResponse) Data(data []byte) error {
	// 赋值Body
	w.setBody(io.NopCloser(bytes.NewReader(data)), int64(len(data)))
	// 发送响应
	return w.send()
}

// JSON 执行响应并序列化响应为JSON
func (w *HttpServerResponse) JSON(obj any) error {
	// 赋值Body
	w.setJSON(obj)
	// 发送响应
	return w.send()
}
