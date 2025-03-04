package httpconn

import (
	"bufio"
	"net/http"
	"time"
)

// 读取HTTP请求（调用结束后请手动关闭Request Body）
func readHttpRequest(ser *HttpServer) (*http.Request, error) {
	// 刷新读超时截止时间
	if err := ser.conn.SetReadDeadline(time.Now().Add(ser.readTimeout)); err != nil {
		return nil, err
	}
	// 准备读取器
	connReader := bufio.NewReader(ser.conn)
	// 是否存在私有协议头
	if ser.priProtoHead != nil {
		// 处理私有协议头
		if err := ser.priProtoHead.ReadRequestHead(connReader); err != nil {
			return nil, err
		}
	}
	// 读取请求
	return http.ReadRequest(connReader)
}

// 发送HTTP响应（res中的响应体将会在发送后自动关闭）
func writeHttpResponse(ser *HttpServer, res *http.Response) error {
	// 刷新写超时截止时间
	if err := ser.conn.SetWriteDeadline(time.Now().Add(ser.writeTimeout)); err != nil {
		res.Body.Close() // 失败需要手动关闭Body
		return err
	}
	// 构建发送器
	connWriter := bufio.NewWriter(ser.conn)
	// 是否存在私有协议头
	if ser.priProtoHead != nil {
		// 处理私有协议头
		if err := ser.priProtoHead.WriteResponseHead(connWriter); err != nil {
			res.Body.Close() // 失败需要手动关闭Body
			return err
		}
	}
	// 发送响应
	err := res.Write(connWriter)
	if err != nil {
		res.Body.Close() // 失败需要手动关闭Body
		return err
	}
	// 刷写数据
	return connWriter.Flush()
}
