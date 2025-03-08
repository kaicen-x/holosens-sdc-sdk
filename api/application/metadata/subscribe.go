/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据订阅
 */
package metadata

// Subscription 订阅信息
type SubscribeInfo struct {
	// 订阅ID
	ID int `json:"id"`
	// 通知接收地址
	//
	//	订阅成功后，SDC将向此地址推送订阅消息
	//	12.0.0版本起支持域名，域名长度最大48个字节，即IP地址或域名入参二选一
	Address string `json:"address"`
	// 通知接收端口
	//
	//	订阅成功后，SDC将向此端口推送订阅消息
	Port uint16 `json:"port"`
	// 订阅剩余时间（必填）
	//
	//	订阅成功后，订阅服务进入倒计时，超时后将不在向目的地址推送事件，
	//	可在超时之前调用此接口刷新订阅时间，单位秒。
	//	取值范围：0~86400，0表示永久有效
	TimeOut int `json:"timeOut"`
	// 是否使用https协议（必填）
	//
	//	由于元数据为敏感数据，只支持通过https协议上报。
	//	此字段固定填1
	HttpsEnable int `json:"httpsEnable"`
	// 服务端开放的元数据接收HTTP/HTTPS服务URL（选填）
	MetaDataURL *string `json:"metaDataURL,omitempty"`
	// 是否上报图片（选填）
	//
	// 不携带此字段或此字段设置为1时上报图片
	// 取值范围：0-不上报图片，1-上报图片
	NeedUploadPic *int `json:"needUploadPic,omitempty"`
}
