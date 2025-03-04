/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备通道信息查询
 */
package devicemanage

// ChannelAttr 设备通道属性
type ChannelAttr struct {
	// 设备通道：
	//	default: 单目设备
	//	details: 二郎神细节路
	//	panorama: 二郎神全景路
	//	magicMaster: 魔方细节路
	//	magicSlave: 魔方全景路
	//	omni: 星图定点
	//	omniMove1: 星图动点1
	//	omniMove2: 星图动点2
	//	1tN: 1tN本地设备视频路
	//	slaveMaster: 子母机的母机视频路
	//	visibleLight: 热成像的可见光路
	//	thermalImaging: 热成像的红外路
	Name string `json:"name"`
	// 设备描述
	//	default: 默认
	//	1tN: 一拖N
	//	ods: 星图、魔方
	Desc string `json:"desc"`
	// 位置
	//	local: 本地，
	//	remote;index1: 远端1
	Location string `json:"location"`
	// 成像类型
	//	visible_light: 可见光
	//	visible_light,focus: 可见光定焦
	//	thermal_imaging: 红外
	ImagingTech string `json:"imaging_tech"`
	// 通道ID：
	//
	//	单CPU架构款型中，通道号按照101、102依次递增；
	//	如二郎神款型细节路为101，全景路为102。
	//	多CPU架构款型中，主芯片通道号按照101、102递
	//	增，从芯片按照2001、2002依次递增；如复眼款型
	//	定点为101，动点1为2001，动点2为2002。一拖N解
	//	决方案中，从机按照10001、10002依次递增。
	ChannelID string `json:"channelId"`
	// 是否需要转发
	//	"1": 需要转发到对应通道
	//	其他：不需要转发。
	ForwardNeed string `json:"forward_need"`
}

// ChannelParam 设备通道参数
type ChannelParam struct {
	Uuid     string            `json:"uuid"`     // 设备通道UUID
	AttrList ChannelAttr       `json:"attrList"` // 设备通道属性（详情请查阅文档）
	FuncList map[string]string `json:"funcList"` // 设备通道功能（详情请查阅文档）
}

// ChannelInfoQueryReply 设备通道信息查询响应
type ChannelInfoQueryReply struct {
	CnsChnParam []ChannelParam `json:"CnsChnParam"` // 设备通道参数列表
}

// ChannelInfoQuery 设备通道信息查询
func (p *Manager) ChannelInfoQuery() (*ChannelInfoQueryReply, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply ChannelInfoQueryReply
	_, err := client.Get("/SDCAPI/V1.0/CnsPaas/ChnQury").
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)

	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
