package main

import (
	"crypto/tls"
	"fmt"
	"log"

	holosenssdcsdk "github.com/bearki/holosens-sdc-sdk"
	"github.com/bearki/holosens-sdc-sdk/api/details/itgt/target/recognize"
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
			instance, err := holosenssdcsdk.NewWithServer(conn)
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

			// 获取设备基础信息
			baseInfo, err := instance.DeviceManager().BaseInfoQuery(101)
			if err != nil {
				log.Printf("BaseInfoQuery error: %s", err)
				return
			}
			fmt.Printf("Keep Live BaseInfo: %+v\n", baseInfo)

			// 获取设备通道信息
			fmt.Println("获取设备通道信息")
			channelInfo, err := instance.DeviceManager().ChannelInfoQuery()
			if err != nil {
				log.Printf("ChannelInfoQuery error: %s", err)
				return
			}
			fmt.Printf("Keep Live ChannelInfo: %+v\n", channelInfo)

			// 查询目标库
			fmt.Println("查询目标库")
			libs, err := instance.ItgtManager().TargetManager().RecognizeManager().TargetLibQuery()
			if err != nil {
				log.Printf("TargetLibsQuery error: %s", err)
				return
			}
			fmt.Printf("Keep Live TargetLibsQuery: %+v\n", libs)

			// 检查目标库是否存在
			fmt.Println("检查目标库是否存在")
			libID := 0
			for _, v := range libs.TargetLibs {
				if v.Name == "test1" {
					libID = v.ID
					break
				}
			}
			if libID == 0 {
				// 库不存在，（创建库）
				fmt.Println("创建目标库")
				err = instance.ItgtManager().TargetManager().RecognizeManager().TargetLibCreate(recognize.TargetLibCreateParams{
					FaceLib: recognize.TargetLibBaseInfo{
						Name: "test1",
					},
				})
				if err != nil {
					log.Printf("TargetLibCreate error: %s", err)
					return
				}
				fmt.Println("创建目标库成功")
			} else {
				// 库存在，（修改库）
				fmt.Println("修改目标库")
				err = instance.ItgtManager().TargetManager().RecognizeManager().TargetLibChange(recognize.TargetLibChangeParams{
					FaceLib: recognize.TargetLibChangeData{
						ID:        libID,
						OnControl: 1,
						TargetLibBaseInfo: recognize.TargetLibBaseInfo{
							Name:      "test1",
							Type:      2,
							Threshold: 83,
							LinkAlarm: 1,
						},
					},
				})
				if err != nil {
					log.Printf("TargetLibChange error: %s", err)
					return
				}
				fmt.Println("修改目标库成功")
			}
		}()
	}
}
