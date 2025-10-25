/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备被动注册（客户端）Socket连接托管
 */
package holosenssdcsdk

import (
	"net"
)

// SessionWithClient 客户端会话
type SessionWithClient struct {
	*Session // 会话
}

// NewWithClient 新建客户端会话
//
//	@param conn: 设备Socket连接通道
//	@return 客户端会话
//	@return 错误信息
func NewWithClient(conn net.Conn) *SessionWithClient {
	// 返回客户端会话
	return &SessionWithClient{
		// 创建会话
		Session: newSession(conn),
	}
}
