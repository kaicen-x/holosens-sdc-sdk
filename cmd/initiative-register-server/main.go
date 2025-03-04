package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	holosenssdcsdk "github.com/bearki/holosens-sdc-sdk"
)

// 主动注册服务端
func main() {
	// 加载证书和私钥
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	// 配置TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// 可以添加其他配置项，如：ClientAuth, MinVersion等
	}

	// 监听TCP端口
	listener, err := tls.Listen("tcp", ":8097", config)
	if err != nil {
		log.Fatalln("server: listen:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on :8097")

	// 开始处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			continue
		}
		// 处理每个连接
		go func() {
			// 构建设备连接实例
			instance, err := holosenssdcsdk.NewDeviceConnectWithServer(conn)
			if err != nil {
				log.Printf("NewDeviceConnect error: %s", err)
				return
			}
			// 打印设备主动注册信息
			fmt.Printf("DeviceInfo: %+v\n", instance.InitiativeRegisterParams)
			// 关闭连接
			defer conn.Close()
			defer fmt.Println("连接断开了")

			// 设置认证信息
			instance.SetAuthorization("ApiAdmin", "a1234567")

			// 保持连接
			for {
				// 获取设备基础信息
				baseInfo, err := instance.DeviceManager().BaseInfoQuery(101)
				if err != nil {
					log.Printf("BaseInfoQuery error: %s", err)
					return
				}
				fmt.Printf("Keep Live BaseInfo: %+v\n", baseInfo)

				time.Sleep(5 * time.Second)

				// 获取设备激活状态
				activateStatusInfo, err := instance.DeviceManager().ActivateStatusQuery()
				if err != nil {
					log.Printf("ActivateStatus error: %s", err)
					return
				}
				fmt.Printf("Keep Live ActivateStatus: %+v\n", activateStatusInfo)

				time.Sleep(5 * time.Second)
			}
		}()
	}
}
