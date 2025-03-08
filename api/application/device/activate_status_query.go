/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口设备激活状态查询
 */
package device

// IdQueryReply 设备激活状态查询响应
type ActivateStatusQueryReply struct {
	// 是否已激活（0-未激活，1-已激活）
	Status int `json:"status"`
}

// ActivateStatusQuery 设备激活状态查询
func (p *Manager) ActivateStatusQuery() (*ActivateStatusQueryReply, error) {
	// 获取Socket连接
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply ActivateStatusQueryReply
	_, err := client.Get("/SDCAPI/V1.0/AuthIaas/ActivaionStatus").
		SetContentType("application/x-www-form-urlencoded").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
