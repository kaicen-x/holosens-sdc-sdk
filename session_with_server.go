/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备主动注册（服务端）Socket连接托管
 */
package holosenssdcsdk

import (
	"net"
	"net/http"

	"github.com/kaicen-x/holosens-sdc-sdk/api/application/device"
	"github.com/kaicen-x/holosens-sdc-sdk/pkg/httpconn"
)

// SessionWithServer 服务端会话
type SessionWithServer struct {
	*Session                                                 // 会话
	InitiativeRegisterParams device.InitiativeRegisterParams // 设备主动注册参数
}

// 设置私有协议头
//
// 在主动注册场景中，由于历史原因，摄像机返回的响应消息最开始有8个字节的私有协
// 议头，解析响应时需对此进行处理，否则可能导致解析出错。
func setPrivateProtocolHead(session *Session) {
	// 获取连接实例上的客户端
	client := session.GetHttp().LockHttpClient()
	defer session.GetHttp().Unlock()
	// 设置HTTP客户端私有协议头
	client.SetPrivateProtocolHead(httpconn.PrivateProtocolHead{
		RequestHead:  nil,
		ResponseHead: make([]byte, 8),
		Strict:       false,
	})
}

// NewWithServer 新建服务端会话（基于TCP服务器）
//
//	@param conn: 设备Socket连接通道
//	@return 服务端会话
//	@return 错误信息
func NewWithServer(conn net.Conn) (*SessionWithServer, error) {
	// 创建会话
	session := newSession(conn)
	// 设置私有协议头
	setPrivateProtocolHead(session)

	// 接收设备主动注册信息
	params, err := session.DeviceManager().InitiativeRegister()
	if err != nil {
		return nil, err
	}

	// 创建并返回服务端会话
	return &SessionWithServer{
		Session:                  session,
		InitiativeRegisterParams: *params,
	}, nil
}

// NewWithHttpServer 新建服务端会话（基于HTTP服务器）
//
//	@param w: 设备HTTP请求响应对象写入器
//	@param r: 设备HTTP请求对象
//	@return 服务端会话
//	@return 错误信息
func NewWithHttpServer(w http.ResponseWriter, r *http.Request) (*SessionWithServer, error) {
	// 处理请求
	params, err := device.InitiativeRegister(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// 接管HTTP的TCP连接
	resController := http.NewResponseController(w)
	conn, _, err := resController.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// 创建会话
	session := newSession(conn)
	// 设置私有协议头
	setPrivateProtocolHead(session)

	// 创建并返回服务端会话
	return &SessionWithServer{
		Session:                  session,
		InitiativeRegisterParams: *params,
	}, nil
}
