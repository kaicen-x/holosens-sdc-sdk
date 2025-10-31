/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备Socket会话缓存
 */
package holosenssdcsdk

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/kaicen-x/holosens-sdc-sdk/api/application/device"
)

var (
	// ErrCacheKeyNotFound：设备会话不存在
	ErrCacheKeyNotFound = errors.New("key not found")
	// ErrCacheInstanceTypeMismatch：设备会话实例类型不匹配
	ErrCacheInstanceTypeMismatch = errors.New("instance type mismatch")
)

// SessionCacheIface 会话缓存接口
type SessionCacheIface interface {
	// 会话必须要实现关闭接口
	Close()
	// 会话必须要实现是否已设置认证信息接口
	IsSetAuthorization() bool
	// 会话必须要实现绑定认证信息修改事件接口
	BindAuthorizationChangeEvent(func(isClear bool))
	// 会话必须要实现设备管理器接口
	DeviceManager() *device.Manager
}

// SessionCacheContext 会话缓存上下文
type SessionCacheContext struct {
	// 唯一标识
	key string
	// 会话实例
	instance SessionCacheIface
	// 心跳上下文取消器互斥锁
	keepliveCancelMtx sync.Mutex
	// 心跳上下文取消器
	keepliveCancel context.CancelFunc
}

// SessionCache 会话缓存器
type SessionCache struct {
	rwMtx    sync.RWMutex                    // 读写锁
	cacheMap map[string]*SessionCacheContext // 缓存数据
}

// NewConnectCache 创建会话缓存器
func NewConnectCache() *SessionCache {
	return &SessionCache{
		cacheMap: make(map[string]*SessionCacheContext),
	}
}

// GetListWithServer 获取全部服务端会话
func (c *SessionCache) GetListWithServer() []*SessionWithServer {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取全部会话
	list := make([]*SessionWithServer, 0, len(c.cacheMap))
	for _, cacheCtx := range c.cacheMap {
		if instance, ok := cacheCtx.instance.(*SessionWithServer); ok {
			list = append(list, instance)
		}
	}
	// 使用KEY进行排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].InitiativeRegisterParams.SerialNumber < list[j].InitiativeRegisterParams.SerialNumber
	})
	// 会话不存在
	return list
}

// GetListWithClient 获取全部客户端会话
func (c *SessionCache) GetListWithClient() []*SessionWithClient {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取全部会话
	list := make([]*SessionWithClient, 0, len(c.cacheMap))
	for _, cacheCtx := range c.cacheMap {
		if instance, ok := cacheCtx.instance.(*SessionWithClient); ok {
			list = append(list, instance)
		}
	}
	// 使用SN进行排序
	// sort.Slice(list, func(i, j int) bool {
	// 	return list[i].InitiativeRegisterParams.SerialNumber < list[j].InitiativeRegisterParams.SerialNumber
	// })
	// 会话不存在
	return list
}

// GetWithServer 获取服务端会话
func (c *SessionCache) GetWithServer(key string) (*SessionWithServer, error) {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取会话
	if cacheCtx, ok := c.cacheMap[key]; ok {
		if instance, ok := cacheCtx.instance.(*SessionWithServer); ok {
			return instance, nil
		}
		return nil, ErrCacheInstanceTypeMismatch
	}
	// 会话不存在
	return nil, ErrCacheKeyNotFound
}

// GetWithClient 获取客户端会话
func (c *SessionCache) GetWithClient(key string) (*SessionWithClient, error) {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取会话
	if cacheCtx, ok := c.cacheMap[key]; ok {
		if instance, ok := cacheCtx.instance.(*SessionWithClient); ok {
			return instance, nil
		}
		return nil, ErrCacheInstanceTypeMismatch
	}
	// 会话不存在
	return nil, ErrCacheKeyNotFound
}

// KeepLive 心跳检测
//
//	心跳失败的设备将会移出缓存
func (c *SessionCache) keeplive(ctx context.Context, key string, instance SessionCacheIface) {
	// 结束时移除会话
	defer c.Remove(key)
	// 心跳失败剩余缓冲次数
	heartbeatMaxFailCount := 3
	// 死循环开始
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
				return
			}
		}
	}
}

// Set 添加会话
func (c *SessionCache) Set(key string, instance SessionCacheIface) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 是否存在同一个设备
	if cacheCtx, ok := c.cacheMap[key]; ok {
		// 移除会话
		c.remove(key, cacheCtx)
	}

	// 创建会话上下文
	cacheCtx := &SessionCacheContext{
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
	// 赋值会话
	c.cacheMap[key] = cacheCtx
	// 1分钟后检查设备是否仍然未配置认证信息
	time.AfterFunc(time.Minute, func() {
		// 是否未配置
		if !instance.IsSetAuthorization() {
			// 移除会话
			c.Remove(key)
		}
	})
}

// 移除会话
func (c *SessionCache) remove(key string, cacheCtx *SessionCacheContext) {
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
	// 移除会话
	delete(c.cacheMap, key)
}

// Delete 移除会话
func (c *SessionCache) Remove(key string) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 移除
	if cacheCtx, ok := c.cacheMap[key]; ok {
		// 执行内部方法移除
		c.remove(key, cacheCtx)
	}
}
