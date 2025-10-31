package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	holosenssdcsdk "github.com/kaicen-x/holosens-sdc-sdk"
	"github.com/kaicen-x/holosens-sdc-sdk/api/application/device"
)

// 运行主动注册服务端
func runSdcServer(socketCache *holosenssdcsdk.SessionCache) {
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
	fmt.Println("Listening on :8097")
	listener, err := tls.Listen("tcp", ":8097", config)
	if err != nil {
		log.Fatalln("server: listen:", err)
	}
	defer listener.Close()

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
			instance, err := holosenssdcsdk.NewWithTcpServer(conn)
			if err != nil {
				log.Printf("NewDeviceConnect error: %s", err)
				return
			}
			// 设置认证信息
			instance.SetAuthorization("ApiAdmin", "a1234567")
			// 打印设备主动注册信息
			fmt.Printf("新的设备注册上来了: %+v\n", instance.InitiativeRegisterParams)

			// 缓存托管实例
			socketCache.Set(instance.InitiativeRegisterParams.SerialNumber, instance)
		}()
	}
}

// 服务端
func main() {
	// 创建连接缓存器
	socketCache := holosenssdcsdk.NewConnectCache()
	// 运行主动注册服务端
	go runSdcServer(socketCache)

	// 运行HTTP服务
	g := gin.Default()
	g.GET("/device/list", func(ctx *gin.Context) {
		// 获取设备连接实例列表
		list := socketCache.GetListWithServer()
		var res []device.InitiativeRegisterParams
		for _, v := range list {
			res = append(res, v.InitiativeRegisterParams)
		}
		ctx.JSON(200, gin.H{
			"devices": res,
		})
	})
	g.Run(":8090")
}
