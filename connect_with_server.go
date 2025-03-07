/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备主动注册（服务端）Socket连接托管
 */
package holosenssdcsdk

import (
	"net"

	"github.com/bearki/holosens-sdc-sdk/api/application/device"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// ConnectTrusteeshipWithServer 服务端Socket连接托管器
type ConnectTrusteeshipWithServer struct {
	*ConnectTrusteeship                                      // 设备Socket连接托管器
	InitiativeRegisterParams device.InitiativeRegisterParams // 设备主动注册参数
}

// 设置私有协议头
//
// 在主动注册场景中，由于历史原因，摄像机返回的响应消息最开始有8个字节的私有协
// 议头，解析响应时需对此进行处理，否则可能导致解析出错。
func setPrivateProtocolHead(connTrusteeship *ConnectTrusteeship) {
	// 获取连接实例上的客户端
	client := connTrusteeship.GetHttp().LockHttpClient()
	defer connTrusteeship.GetHttp().Unlock()
	// 设置HTTP客户端私有协议头
	client.SetPrivateProtocolHead(httpconn.PrivateProtocolHead{
		RequestHead:  nil,
		ResponseHead: make([]byte, 8),
		Strict:       false,
	})
}

// NewWithServer 新建服务端Socket连接托管器
//
//	@param conn: 设备Socket连接通道
//	@return 设备Socket连接托管器
//	@return 错误信息
func NewWithServer(conn net.Conn) (*ConnectTrusteeshipWithServer, error) {
	// 创建设备Socket连接托管器
	connTrusteeship := newConnectTrusteeship(conn)
	// 设置私有协议头
	setPrivateProtocolHead(connTrusteeship)

	// 接收设备主动注册信息
	params, err := connTrusteeship.DeviceManager().InitiativeRegister()
	if err != nil {
		return nil, err
	}

	// 返回托管器
	return &ConnectTrusteeshipWithServer{
		ConnectTrusteeship:       connTrusteeship,
		InitiativeRegisterParams: *params,
	}, nil
}
