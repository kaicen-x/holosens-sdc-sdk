package httpconn

import (
	"bufio"
	"context"
	"net/http"
	"time"

	"github.com/bearki/holosens-sdc-sdk/pkg/digest"
)

// 发送HTTP请求（req中的请求体将会在发送后自动关闭）
func writeHttpRequest(connInfo *ConnectInfo, req *http.Request) error {
	// 先获取WWW-Authenticate
	{
		// 刷新写超时截止时间
		if err := connInfo.conn.SetWriteDeadline(time.Now().Add(connInfo.connWriteTimeout)); err != nil {
			req.Body.Close()
			return err
		}
		// 发送空请求获取WWW-Authenticate
		tmpReq := req.Clone(context.Background())
		tmpReq.Body = http.NoBody
		tmpReq.ContentLength = 0
		// 发送请求
		connWriter := bufio.NewWriter(connInfo.conn)
		err := tmpReq.Write(connWriter)
		if err != nil {
			req.Body.Close()
			return err
		}
		// 刷写数据
		err = connWriter.Flush()
		if err != nil {
			req.Body.Close()
			return err
		}
		// 读取响应
		res, err := readHttpResponse(connInfo, req)
		if err != nil {
			req.Body.Close()
			return err
		}
		// 解析响应
		realm, nonce, algorithm := digest.ParseDigestWwwAuthenticate(res.Header)
		// 构建认证并赋值到原始请求
		req.Header.Set("Authorization", digest.MakeDigestAuthorization(
			req.Method, req.URL.RequestURI(),
			realm, nonce, algorithm,
			connInfo.username, connInfo.password,
		))
	}

	// 获取认证信息后再发送请求
	{
		// 刷新写超时截止时间
		if err := connInfo.conn.SetWriteDeadline(time.Now().Add(connInfo.connWriteTimeout)); err != nil {
			req.Body.Close()
			return err
		}
		// 发送请求
		connWriter := bufio.NewWriter(connInfo.conn)
		err := req.Write(connWriter)
		if err != nil {
			req.Body.Close()
			return err
		}
		// 刷写数据
		return connWriter.Flush()
	}
}

// 发送HTTP响应（res中的响应体将会在发送后自动关闭）
func writeHttpResponse(connInfo *ConnectInfo, res *http.Response) error {
	// 刷新写超时截止时间
	if err := connInfo.conn.SetWriteDeadline(time.Now().Add(connInfo.connWriteTimeout)); err != nil {
		res.Body.Close()
		return err
	}
	// 发送响应
	connWriter := bufio.NewWriter(connInfo.conn)
	err := res.Write(connWriter)
	if err != nil {
		res.Body.Close()
		return err
	}
	// 刷写数据
	return connWriter.Flush()
}
