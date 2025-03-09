/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备Socket连接托管缓存
 */
package holosenssdcsdk

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/bearki/holosens-sdc-sdk/api/application/device"
)

var (
	// ErrCacheKeyNotFound：设备托管器不存在
	ErrCacheKeyNotFound = errors.New("key not found")
	// ErrCacheInstanceTypeMismatch：设备托管器实例类型不匹配
	ErrCacheInstanceTypeMismatch = errors.New("instance type mismatch")
)

// ConnectTrusteeshipWithServer 设备Socket连接托管器实例接口
type ConnectCacheInstance interface {
	// 托管器必须要实现关闭接口
	Close()
	// 托管器必须要实现是否已设置认证信息接口
	IsSetAuthorization() bool
	// 托管器必须要实现绑定认证信息修改事件接口
	BindAuthorizationChangeEvent(func(isClear bool))
	// 托管器必须要实现设备管理器接口
	DeviceManager() *device.Manager
}

// ConnectCacheContext 设备Socket连接托管缓存器上下文
type ConnectCacheContext struct {
	// 唯一标识
	key string
	// 设备Socket连接托管器实例
	instance ConnectCacheInstance
	// 心跳上下文取消器互斥锁
	keepliveCancelMtx sync.Mutex
	// 心跳上下文取消器
	keepliveCancel context.CancelFunc
}

// ConnectCache 设备Socket连接托管缓存器
type ConnectCache struct {
	rwMtx    sync.RWMutex                    // 读写锁
	cacheMap map[string]*ConnectCacheContext // 缓存器
}

// NewConnectCache 创建设备Socket连接托管缓存器
func NewConnectCache() *ConnectCache {
	return &ConnectCache{
		cacheMap: make(map[string]*ConnectCacheContext),
	}
}

// GetWithServer 服务端获取全部托管器
func (c *ConnectCache) GetListWithServer() []*ConnectTrusteeshipWithServer {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取全部托管器
	list := make([]*ConnectTrusteeshipWithServer, 0, len(c.cacheMap))
	for _, cacheCtx := range c.cacheMap {
		if instance, ok := cacheCtx.instance.(*ConnectTrusteeshipWithServer); ok {
			list = append(list, instance)
		}
	}
	// 使用KEY进行排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].InitiativeRegisterParams.SerialNumber < list[j].InitiativeRegisterParams.SerialNumber
	})
	// 托管器不存在
	return list
}

// GetWithClient 客户端获取全部托管器
func (c *ConnectCache) GetListWithClient() []*ConnectTrusteeshipWithClient {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取全部托管器
	list := make([]*ConnectTrusteeshipWithClient, 0, len(c.cacheMap))
	for _, cacheCtx := range c.cacheMap {
		if instance, ok := cacheCtx.instance.(*ConnectTrusteeshipWithClient); ok {
			list = append(list, instance)
		}
	}
	// 使用SN进行排序
	// sort.Slice(list, func(i, j int) bool {
	// 	return list[i].InitiativeRegisterParams.SerialNumber < list[j].InitiativeRegisterParams.SerialNumber
	// })
	// 托管器不存在
	return list
}

// GetWithServer 服务端获取托管器
func (c *ConnectCache) GetWithServer(key string) (*ConnectTrusteeshipWithServer, error) {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取托管器
	if cacheCtx, ok := c.cacheMap[key]; ok {
		if instance, ok := cacheCtx.instance.(*ConnectTrusteeshipWithServer); ok {
			return instance, nil
		}
		return nil, ErrCacheInstanceTypeMismatch
	}
	// 托管器不存在
	return nil, ErrCacheKeyNotFound
}

// GetWithClient 客户端获取托管器
func (c *ConnectCache) GetWithClient(key string) (*ConnectTrusteeshipWithClient, error) {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取托管器
	if cacheCtx, ok := c.cacheMap[key]; ok {
		if instance, ok := cacheCtx.instance.(*ConnectTrusteeshipWithClient); ok {
			return instance, nil
		}
		return nil, ErrCacheInstanceTypeMismatch
	}
	// 托管器不存在
	return nil, ErrCacheKeyNotFound
}

// KeepLive 心跳检测
//
//	心跳失败的设备将会移出缓存
func (c *ConnectCache) keeplive(ctx context.Context, key string, instance ConnectCacheInstance) {
	// 结束时移除托管器
	defer c.Remove(key)
	// 心跳失败剩余缓冲次数
	heartbeatMaxFailCount := 3
	// 死循环开始
KEEPLIVE_LOOP:
	for {
		select {
		// 检查上下文是否终止
		case <-ctx.Done():
			// 结束
			return

		// 是否达到检测间隔时长
		case <-time.After(time.Minute):
			// 心跳检测
			_, err := instance.DeviceManager().BaseInfoQuery(101)
			if err != nil {
				// 心跳失败，是否还有剩余缓冲次数
				if heartbeatMaxFailCount > 0 {
					// 减1
					heartbeatMaxFailCount--
					// 继续检测
					continue
				}
				// 跳出循环
				break KEEPLIVE_LOOP
			}
		}
	}
}

// Set 添加托管器
func (c *ConnectCache) Set(key string, instance ConnectCacheInstance) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 是否存在同一个设备
	if cacheCtx, ok := c.cacheMap[key]; ok {
		// 移除托管器
		c.remove(key, cacheCtx)
	}

	// 创建托管器上下文
	cacheCtx := &ConnectCacheContext{
		key:            key,
		instance:       instance,
		keepliveCancel: nil,
	}
	// 是否已存在认证信息
	if instance.IsSetAuthorization() {
		// 创建心跳检测上下文
		keepliveCtx, keepliveCancel := context.WithCancel(context.Background())
		cacheCtx.keepliveCancel = keepliveCancel
		// 启动心跳检测
		go c.keeplive(keepliveCtx, key, instance)
	}
	// 监听认证信息修改事件（回调事件）
	instance.BindAuthorizationChangeEvent(func(isClear bool) {
		// 心跳上下文取消器互斥锁
		cacheCtx.keepliveCancelMtx.Lock()
		defer cacheCtx.keepliveCancelMtx.Unlock()
		// 停止心跳检测
		if cacheCtx.keepliveCancel != nil {
			cacheCtx.keepliveCancel()
			cacheCtx.keepliveCancel = nil
		}
		// 非清空认证信息的情况 并且 实例还有效时
		// 需要启动新的心跳检测
		if !isClear && cacheCtx.instance != nil {
			// 创建心跳检测上下文
			keepliveCtx, keepliveCancel := context.WithCancel(context.Background())
			cacheCtx.keepliveCancel = keepliveCancel
			// 启动心跳检测
			go c.keeplive(keepliveCtx, key, cacheCtx.instance)
		}
	})
	// 赋值托管器
	c.cacheMap[key] = cacheCtx
}

// Delete 移除托管器
func (c *ConnectCache) remove(key string, cacheCtx *ConnectCacheContext) {
	// 心跳上下文取消器互斥锁
	cacheCtx.keepliveCancelMtx.Lock()
	defer cacheCtx.keepliveCancelMtx.Unlock()
	// 停止心跳检测
	if cacheCtx.keepliveCancel != nil {
		cacheCtx.keepliveCancel()
		cacheCtx.keepliveCancel = nil
	}
	// 关闭连接
	cacheCtx.instance.Close()
	// 清空实例
	cacheCtx.instance = nil
	// 移除托管器
	delete(c.cacheMap, key)
}

// Delete 移除托管器
func (c *ConnectCache) Remove(key string) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 移除
	if cacheCtx, ok := c.cacheMap[key]; ok {
		// 执行内部方法移除
		c.remove(key, cacheCtx)
	}
}
