/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备连接托管
 */
package holosenssdcsdk

import (
	devicemanage "github.com/bearki/holosens-sdc-sdk/api/device-manage"
	intelligentmetadata "github.com/bearki/holosens-sdc-sdk/api/intelligent-metadata"
	"github.com/bearki/holosens-sdc-sdk/pkg/httpconn"
)

// DeviceConnect 设备Socket连接托管基础
type DeviceConnect struct {
	connInstance    *httpconn.Connect            // Socket连接实例
	deviceManager   *devicemanage.Manager        // 设备管理与维护管理器
	metadataManager *intelligentmetadata.Manager // 智能元数据对接管理器
}

// SetAuthorization 设置连接认证信息
func (p *DeviceConnect) SetAuthorization(username, password string) {
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()
	client.SetDigestAuth(username, password)
}

// DeviceManager 获取设备管理与维护管理器
func (p *DeviceConnect) DeviceManager() *devicemanage.Manager {
	return p.deviceManager
}

// IntelligentMetadataManager 获取智能元数据对接管理器
func (p *DeviceConnect) IntelligentMetadataManager() *intelligentmetadata.Manager {
	return p.metadataManager
}
