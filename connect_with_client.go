/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备被动注册（客户端）Socket连接托管
 */
package holosenssdcsdk

import (
	"net"
)

// ConnectTrusteeshipWithClient 客户端Socket连接托管器
type ConnectTrusteeshipWithClient struct {
	*ConnectTrusteeship // 设备Socket连接托管器
}

// NewWithClient 新建客户端Socket连接托管器
//
//	@param conn: 设备Socket连接通道
//	@return 设备Socket连接托管器
//	@return 错误信息
func NewWithClient(conn net.Conn) *ConnectTrusteeshipWithClient {
	// 返回托管器
	return &ConnectTrusteeshipWithClient{
		ConnectTrusteeship: newConnectTrusteeship(conn),
	}
}
