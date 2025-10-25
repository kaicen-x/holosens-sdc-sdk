/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口目标记录模块
 */
package recognize

// TargetRecordParamsWithLib 目标记录参数关联的目标库信息
type TargetRecordParamsWithLib struct {
	// 目标库ID（必填）
	ID int `json:"faceLibId"`
}

// TargetRecordBaseInfo 目标记录基础信息
type TargetRecordBaseInfo struct {
	// 姓名（选填）
	//	取值范围：0~63字符
	Name string `json:"name"`
	// 性别（选填）
	//	取值范围：0-男，1-女，2-未知，-1-未知
	Gender int `json:"gender"`
	// 生日（选填）
	//	取值范围：0~31字符，支持数字、字母、"-"以及"/"
	Birthday string `json:"birthday"`
	// 省级（选填）
	//	取值范围：0~31字符
	Province string `json:"province"`
	// 地级（选填）
	//	取值范围：0~47字符
	City string `json:"city"`
	// 证件类型（选填）
	//	取值范围：0-身份证，1-护照，2-军官证，3-驾驶证，4-其他，-1-未知
	CardType int `json:"cardType"`
	// 证件号（选填）
	//	取值范围：0~31字符
	CardId string `json:"cardId"`
	// 外部指定的目标ID（选填）
	//	用于外部管理目标库记录
	//	SDC 9.0.0-LG0001版本新增
	//	取值范围：不超过64个字符，支持英文字母、数字
	ExternalId string `json:"externalId"`
	// 目标图片是否保存到摄像机（选填）
	//	目标图片是否保存到摄像机，开启时关闭时，只保存目标特征值及目标信息
	//	SDC 9.0.0-LG0001版本新增
	//	取值范围：0-不保存，1-保存，-1-未知
	IsStore int `json:"isStore"`
}

// TargetRecordCreateInfo 目标记录创建信息
type TargetRecordCreateInfo struct {
	// 图片名（必填）
	//	取值范围：0~63字符
	PicName string `json:"picName"`
	// 目标记录基础信息
	TargetRecordBaseInfo
}

// TargetRecordChangeInfo 目标记录修改信息
type TargetRecordChangeInfo struct {
	// 目标记录ID（必填）
	ID int `json:"faceId"`
	// 目标记录创建信息
	TargetRecordCreateInfo
}

// TargetRecordBatchQueryInfo 目标记录批量查询信息
type TargetRecordBatchQueryInfo struct {
	// 目标记录基本信息（必填）
	TargetRecordBaseInfo
	// 特征值状态（选填）
	//	取值范围：
	// 	0-新建目标库状态
	// 	1-查询全部
	// 	2-重新提取
	// 	3-未提取
	// 	4-已提取
	// 	5-提取失败
	// 	6-图片尺寸不规范
	// 	7-图片解码失败
	// 	8-目标检测失败
	// 	9-特征提取失败
	// 	10-特征写入失败
	// 	11-目标图片清晰度不够
	// 	12-目标图片遮挡较严重
	FeatureStatus int `json:"featureStatus"`
}

// TargetRecordInfo 目标记录信息
type TargetRecordInfo struct {
	// 目标记录ID
	ID int `json:"id"`
	// 图片名（必填）
	//	取值范围：0~63字符
	PicName string `json:"picName"`
	// 目标记录基础信息
	TargetRecordBaseInfo
	// 特征值状态
	//	取值范围：
	// 	0-新建目标库状态
	// 	1-查询全部
	// 	2-重新提取
	// 	3-未提取
	// 	4-已提取
	// 	5-提取失败
	// 	6-图片尺寸不规范
	// 	7-图片解码失败
	// 	8-目标检测失败
	// 	9-特征提取失败
	// 	10-特征写入失败
	// 	11-目标图片清晰度不够
	// 	12-目标图片遮挡较严重
	FeatureStatus int `json:"featureStatus"`
}
