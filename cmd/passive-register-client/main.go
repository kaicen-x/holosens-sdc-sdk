package main

import (
	"fmt"
	"net"

	holosenssdcsdk "github.com/kaicen-x/holosens-sdc-sdk"
)

// 被动注册客户端
func main() {
	// 建立TLS TCP连接
	conn, err := net.Dial("tcp", "192.168.8.27:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 托管连接
	deviceConnect := holosenssdcsdk.NewWithTcpClient(conn)
	// 设置北向接口认证信息
	deviceConnect.SetAuthorization("ApiAdmin", "a1234567")

	// 获取设备管理与维护管理器
	deviceManager := deviceConnect.DeviceManager()
	// 获取设备基础信息
	baseInfo, err := deviceManager.BaseInfoQuery(101)
	if err != nil {
		panic(err)
	}

	fmt.Printf("BaseInfo: %+v\n", baseInfo)
}
