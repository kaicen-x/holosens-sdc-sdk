package httpconn

import (
	"bufio"
	"context"
	"net/http"
	"time"

	"github.com/bearki/holosens-sdc-sdk/pkg/digest"
)

// 内部发送HTTP请求（req中的请求体将会在发送后自动关闭）
func internalWriteHttpRequest(cli *HttpClient, req *http.Request) error {
	// 刷新写超时截止时间
	if err := cli.conn.SetWriteDeadline(time.Now().Add(cli.writeTimeout)); err != nil {
		req.Body.Close() // 发送前失败需要手动关闭传入的请求体
		return err
	}
	// 构建发送器
	connWriter := bufio.NewWriter(cli.conn)
	// 是否存在私有协议头
	if cli.priProtoHead != nil {
		// 处理私有协议头
		if err := cli.priProtoHead.WriteRequestHead(connWriter); err != nil {
			req.Body.Close() // 发送前失败需要手动关闭传入的请求体
			return err
		}
	}
	// 发送请求
	err := req.Write(connWriter)
	if err != nil {
		req.Body.Close() // 发送前失败需要手动关闭传入的请求体
		return err
	}
	// 刷写数据
	return connWriter.Flush()
}

// 发送HTTP请求（req中的请求体将会在发送后自动关闭）
func writeHttpRequest(cli *HttpClient, req *http.Request) error {
	// 是否需要 Digest 认证
	if cli.auth != nil && cli.auth.Type == HttpClientAuthTypeDigest {
		/**
		 *先获取WWW-Authenticate
		 */
		// 拷贝Request将Body置空
		tmpReq := req.Clone(context.Background())
		tmpReq.Body = http.NoBody
		tmpReq.ContentLength = 0
		// 内部发送请求
		err := internalWriteHttpRequest(cli, tmpReq)
		if err != nil {
			req.Body.Close() // 发送前失败需要手动关闭传入的请求体
			return err
		}
		// 读取响应
		res, err := readHttpResponse(cli, tmpReq)
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
			cli.auth.Username, cli.auth.Password,
		))
	}

	// 内部发送请求
	return internalWriteHttpRequest(cli, req)
}

// 读取HTTP响应（调用结束后请手动关闭Response Body）
func readHttpResponse(cli *HttpClient, req *http.Request) (*http.Response, error) {
	// 刷新读超时截止时间
	if err := cli.conn.SetReadDeadline(time.Now().Add(cli.readTimeout)); err != nil {
		return nil, err
	}
	// 准备读取器
	connReader := bufio.NewReader(cli.conn)
	// 是否存在私有协议头
	if cli.priProtoHead != nil {
		// 处理私有协议头
		if err := cli.priProtoHead.ReadResponseHead(connReader); err != nil {
			return nil, err
		}
	}
	// 读取响应
	return http.ReadResponse(connReader, req)
}
