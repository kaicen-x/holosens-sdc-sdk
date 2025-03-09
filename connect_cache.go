/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备Socket连接托管缓存
 */
package holosenssdcsdk

import (
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
	// 托管器必须要实现设备管理器接口
	DeviceManager() *device.Manager
}

// ConnectCache 设备Socket连接托管缓存器
type ConnectCache struct {
	rwMtx    sync.RWMutex                    // 读写锁
	cacheMap map[string]ConnectCacheInstance // 缓存器
}

// NewConnectCache 创建设备Socket连接托管缓存器
func NewConnectCache() *ConnectCache {
	return &ConnectCache{
		cacheMap: make(map[string]ConnectCacheInstance),
	}
}

// GetWithServer 服务端获取全部托管器
func (c *ConnectCache) GetListWithServer() []*ConnectTrusteeshipWithServer {
	// 加读锁
	c.rwMtx.RLock()
	defer c.rwMtx.RUnlock()
	// 获取全部托管器
	list := make([]*ConnectTrusteeshipWithServer, 0, len(c.cacheMap))
	for _, val := range c.cacheMap {
		if instance, ok := val.(*ConnectTrusteeshipWithServer); ok {
			list = append(list, instance)
		}
	}
	// 使用SN进行排序
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
	for _, val := range c.cacheMap {
		if instance, ok := val.(*ConnectTrusteeshipWithClient); ok {
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
	if val, ok := c.cacheMap[key]; ok {
		if instance, ok := val.(*ConnectTrusteeshipWithServer); ok {
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
	if val, ok := c.cacheMap[key]; ok {
		if instance, ok := val.(*ConnectTrusteeshipWithClient); ok {
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
func (c *ConnectCache) keeplive(key string, val ConnectCacheInstance) {
	// 结束时移除托管器
	defer c.Remove(key)
	// 心跳失败剩余缓冲次数
	heartbeatMaxFailCount := 3
	// 死循环开始
	for {
		// 需要一点间隔时长
		<-time.After(time.Minute)
		// 心跳检测
		_, err := val.DeviceManager().BaseInfoQuery(101)
		if err != nil {
			// 心跳失败，是否还有剩余缓冲次数
			if heartbeatMaxFailCount > 0 {
				// 减1
				heartbeatMaxFailCount--
				// 继续检测
				continue
			}
			// 跳出循环
			break
		}
	}
}

// Set 添加托管器
func (c *ConnectCache) Set(key string, val ConnectCacheInstance) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 设置托管器
	c.cacheMap[key] = val
	// 启动心跳检测
	go c.keeplive(key, val)
}

// Delete 移除托管器
func (c *ConnectCache) Remove(key string) {
	// 加写锁
	c.rwMtx.Lock()
	defer c.rwMtx.Unlock()
	// 管理托管器连接
	if val, ok := c.cacheMap[key]; ok {
		// 关闭连接
		val.Close()
		// 移除托管器
		delete(c.cacheMap, key)
	}
}
