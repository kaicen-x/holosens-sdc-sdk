/**
 * @Author: Bearki
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标记录批量查询
 */
package recognize

// TargetRecordBatchQueryParams 目标记录批量查询参数
//
//	1、全部查询，gender=-1，cardType=-1，isStore=-1，其他字段，数字为0，字符串的为空
//	2、条件查询，未填的字段：除以上三个字段外，其他类型数字为0，字符串的为空
type TargetRecordBatchQueryParams struct {
	// 目标库信息（必填）
	Lib TargetRecordParamsWithLib `json:"faceLib"`
	// 目标记录查询索引（选填）
	StartIndex int `json:"startIndex"`
	// 目标记录一次查询条数（选填）
	// 	取值范围：[0,5000]
	FindNum int `json:"findNum"`
	// 目标记录批量查询信息（必填）
	Condition TargetRecordBatchQueryInfo `json:"condition"`
}

// TargetRecordBatchQueryReply 目标记录批量查询响应
type TargetRecordBatchQueryReply struct {
	// 目标记录总数
	TotalNum int `json:"totalNum"`
	// 目标记录列表
	TargetRecordList []TargetRecordInfo `json:"faceRecordArry"`
}

// TargetRecordBatchQuery 目标记录批量查询
//
//	1、全部查询，gender=-1，cardType=-1，isStore=-1，其他字段，数字为0，字符串的为空
//	2、条件查询，未填的字段：除以上三个字段外，其他类型数字为0，字符串的为空
//	@param	params: 目标记录批量查询参数
//	@return	错误信息
func (p *Manager) TargetRecordBatchQuery(params TargetRecordBatchQueryParams) (*TargetRecordBatchQueryReply, error) {
	// 获取客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 发送请求
	var reply TargetRecordBatchQueryReply
	_, err := client.Post("/SDCAPI/V1.0/FaceApp/FaceRecog/FaceLibs/FaceRecordQuery").
		SetJSON(&params).
		SetContentType("application/json").
		DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
