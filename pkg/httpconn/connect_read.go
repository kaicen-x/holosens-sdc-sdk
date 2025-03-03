package httpconn

import (
	"bufio"
	"errors"
	"net/http"
	"time"
)

// 读取HTTP请求（调用结束后请手动关闭Request Body）
func readHttpRequest(connInfo *ConnectInfo) (*http.Request, error) {
	// 刷新读超时截止时间
	if err := connInfo.conn.SetReadDeadline(time.Now().Add(connInfo.connReadTimeout)); err != nil {
		return nil, err
	}
	// 读取请求
	connReader := bufio.NewReader(connInfo.conn)
	return http.ReadRequest(connReader)
}

// 读取HTTP响应（调用结束后请手动关闭Response Body）
func readHttpResponse(connInfo *ConnectInfo, req *http.Request) (*http.Response, error) {
	// 刷新读超时截止时间
	if err := connInfo.conn.SetReadDeadline(time.Now().Add(connInfo.connReadTimeout)); err != nil {
		return nil, err
	}
	// 读取响应
	connReader := bufio.NewReader(connInfo.conn)
	// 是否需要移除私有协议头
	if connInfo.connType {
		privateProtocol := make([]byte, 8)
		n, err := connReader.Read(privateProtocol)
		if err != nil {
			return nil, err
		}
		if n != 8 {
			return nil, errors.New("private protocol error")
		}
	}
	// 读取响应
	return http.ReadResponse(connReader, req)
}
