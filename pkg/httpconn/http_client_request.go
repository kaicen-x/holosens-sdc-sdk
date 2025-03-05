package httpconn

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
)

// HttpClientRequest HTTP自定义客户端请求对象
type HttpClientRequest struct {
	cli   *HttpClient   // 客户端
	req   *http.Request // HTTP原生请求对象
	query url.Values    // Query参数临时缓存一下（有助于提高性能）
	err   error         // 缓存链式调用中间产生的异常
}

// NewHttpClientRequest 创建HTTP自定义客户端请求对象
func NewHttpClientRequest(client *HttpClient, method, url string) *HttpClientRequest {
	// 构建请求体
	req, err := http.NewRequest(method, url, nil)
	// 构建HTTP自定义客户端请求对象
	return &HttpClientRequest{
		cli:   client,
		err:   err,
		req:   req,
		query: nil,
	}
}

// SetHeader 设置请求头
func (r *HttpClientRequest) SetHeader(key, val string) *HttpClientRequest {
	// 检查
	if r.req == nil {
		return r
	}
	// 设置请求头
	if r.req.Header == nil {
		r.req.Header = make(http.Header)
	}
	r.req.Header.Set(key, val)
	// OK
	return r
}

// AddHeader 添加请求头
func (r *HttpClientRequest) AddHeader(key, val string) *HttpClientRequest {
	// 检查
	if r.req == nil {
		return r
	}
	// 添加请求头
	if r.req.Header == nil {
		r.req.Header = make(http.Header)
	}
	r.req.Header.Add(key, val)
	// OK
	return r
}

// SetQuery 设置Query请求参数
func (r *HttpClientRequest) SetQuery(key, val string) *HttpClientRequest {
	// 设置参数
	if r.query == nil {
		r.query = make(url.Values)
	}
	r.query.Set(key, val)
	// OK
	return r
}

// AddQuery 添加Query请求参数
func (r *HttpClientRequest) AddQuery(key, val string) *HttpClientRequest {
	// 添加参数
	if r.query == nil {
		r.query = make(url.Values)
	}
	r.query.Add(key, val)
	// OK
	return r
}

// SetContentType 设置Content-Type请求头信息
func (r *HttpClientRequest) SetContentType(contentType string) *HttpClientRequest {
	// 检查
	if r.req == nil {
		return r
	}
	// 设置请求头信息
	if r.req.Header == nil {
		r.req.Header = make(http.Header)
	}
	r.req.Header.Set("Content-Type", contentType)
	// OK
	return r
}

// SetBody 设置请求体
func (r *HttpClientRequest) SetBody(body io.ReadCloser, bodySize int64) *HttpClientRequest {
	// 检查
	if r.req == nil {
		body.Close() // 提前结束需要关闭请求体
		return r
	}
	// 设置请求体
	r.req.Body = body
	r.req.ContentLength = bodySize
	// OK
	return r
}

// SetJsonBody 设置JSON请求体
func (r *HttpClientRequest) SetJSON(data any) *HttpClientRequest {
	// 检查
	if r.req == nil {
		return r
	}
	// json编码数据
	dataBytes, err := json.Marshal(data)
	if err != nil {
		r.err = errors.Join(r.err, err)
		return r
	}
	// 设置请求体
	r.SetBody(io.NopCloser(bytes.NewReader(dataBytes)), int64(len(dataBytes)))
	// 设置请求头信息
	r.SetContentType("application/json; charset=UTF-8")
	// OK
	return r
}

// Send 执行请求（请主动关闭Response.Body）
func (r *HttpClientRequest) Send() (*http.Response, error) {
	// 提取原生请求对象
	req := r.req
	if req == nil {
		return nil, errors.Join(r.err, errors.New("request is nil"))
	}
	// 检查是否存在缓存错误
	if r.err != nil {
		if req.Body != nil {
			req.Body.Close() // 提前结束需要手动释放Body
		}
		return nil, r.err
	}
	// 检查URL是否有效
	if req.URL == nil {
		if req.Body != nil {
			req.Body.Close() // 提前结束需要手动释放Body
		}
		return nil, errors.New("invalid URL")
	}
	// 是否存在Query参数需要赋值
	if len(r.query) > 0 {
		req.URL.RawQuery = r.query.Encode()
	}
	// 发送请求
	err := writeHttpRequest(r.cli, req)
	if err != nil {
		return nil, err
	}
	// 读取响应
	res, err := readHttpResponse(r.cli, req)
	if err != nil {
		return nil, err
	}
	// OK
	return res, nil
}

// DecodeJSON 执行请求并解析响应为JSON（请不要重复关闭Response.Body）
func (r *HttpClientRequest) DecodeJSON(obj any) (*http.Response, error) {
	// 执行请求
	res, err := r.Send()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	// 检查响应状态码
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if len(body) > 0 {
			return res, errors.New(string(body))
		}
		return res, errors.New(res.Status)
	}
	// 解析响应体
	return res, json.Unmarshal(body, obj)
}

// DecodeFormData 执行请求并解析响应为FormData（请不要重复关闭Response.Body）
//
//	@param maxMemorySize: 最大内存大小，单位为字节（建议不要超过20MB，超出为文件内容将会自动落盘到临时文件，不用担心文件接收不到）
func (r *HttpClientRequest) DecodeFormData(maxMemorySize int64) (*multipart.Form, *http.Response, error) {
	// 执行请求
	res, err := r.Send()
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	// 解析boundary
	mediatype, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, nil, err
	}
	if mediatype != "multipart/form-data" {
		return nil, nil, errors.New("invalid Content-Type")
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, nil, errors.New("invalid boundary") // 无法解析boundary
	}
	// 创建表单读取器
	reader := multipart.NewReader(res.Body, boundary)
	form, err := reader.ReadForm(maxMemorySize)
	if err != nil {
		return nil, nil, err
	}
	// OK
	return form, res, nil
}
