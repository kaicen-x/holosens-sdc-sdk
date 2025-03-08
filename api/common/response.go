package common

import "net/http"

// ResponseStatus 通用响应状态对象
type ResponseStatus struct {
	RequestURL   string `json:"RequestURL,omitempty"` // 此响应对应请求的URL
	StatusCode   int    `json:"StatusCode"`           // 响应码
	StatusString string `json:"StatusString"`         // 请求出错时返回的描述信息
}

// Response 通用响应对象
type Response[T any] struct {
	ResponseStatus T `json:"ResponseStatus"` // 通用响应对象
}

// NewResponseWithSuccess 创建一个成功的响应对象
func NewResponseWithSuccess(req *http.Request) *Response[ResponseStatus] {
	uri := ""
	if req != nil {
		uri = req.URL.RequestURI()
	}
	return &Response[ResponseStatus]{
		ResponseStatus: ResponseStatus{
			RequestURL:   uri,
			StatusCode:   0,
			StatusString: "OK",
		},
	}
}

// NewResponseWithFailed 创建一个失败的响应对象
func NewResponseWithFailed(req *http.Request) *Response[ResponseStatus] {
	uri := ""
	if req != nil {
		uri = req.URL.String()
	}
	return &Response[ResponseStatus]{
		ResponseStatus: ResponseStatus{
			RequestURL:   uri,
			StatusCode:   -1,
			StatusString: "FAILED",
		},
	}
}
