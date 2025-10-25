package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	holosenssdcsdk "github.com/kaicen-x/holosens-sdc-sdk"
)

// 服务端
func main() {
	// 创建连接缓存器
	socketCache := holosenssdcsdk.NewConnectCache()
	// 初始化HTTP服务
	g := gin.Default()
	g.PUT("/register", func(ctx *gin.Context) {
		// 接管设备注册，并处理底层Socket连接
		instance, err := holosenssdcsdk.NewWithHttpServer(ctx.Writer, ctx.Request)
		if err != nil {
			return
		}
		// 设置认证信息
		instance.SetAuthorization("ApiAdmin", "a1234567")
		// 打印设备主动注册信息
		fmt.Printf("新的设备注册上来了: %+v\n", instance.InitiativeRegisterParams)
		// 缓存托管实例
		socketCache.Set(instance.InitiativeRegisterParams.SerialNumber, instance)
	})
	// 启动HTTP服务
	g.RunTLS(":8090", "server.crt", "server.key")
}
