package recognize

// TargetLibBaseInfo 目标库基础信息
type TargetLibBaseInfo struct {
	// 目标库名（必填）
	//	取值范围：0~63个字符
	Name string `json:"libName"`
	// 目标库操作类型（选填，默认：1）
	//	取值范围：1-黑名单, 2-白名单
	Type int `json:"libType"`
	// 目标库设防阈值（选填，默认：0）
	// 	取值范围：[0,99]
	Threshold int `json:"threshold"`
	// 是否输出告警（选填，默认：0）
	//	取值范围：0-不输出告警, 1-输出告警
	LinkAlarm int `json:"libLinkAlarm"`
}

// TargetLibInfo 目标库信息
type TargetLibInfo struct {
	// 目标库ID
	ID int `json:"ID"`
	// 目标库名
	//	取值范围：0~63个字符
	Name string `json:"FaceListName"`
	// 目标库操作类型
	//	取值范围：1-黑名单, 2-白名单
	Type int `json:"FaceListType"`
	// 目标库设防状态
	OnControl int `json:"OnControl"`
	// 目标库设防阈值
	// 	取值范围：[0,99]
	Threshold int `json:"Threshold"`
	// 算法版本
	AlgVersion string `json:"AlgVersion"`
	// 特征值状态
	FeaStatus int `json:"FeaStatus"`
	// 已提取数量
	ExtractedNum int `json:"ExtractedNum"`
	// 图片总数量
	TotalPicNum int `json:"TotalPicNum"`
	// 是否输出告警
	//	取值范围：0-不输出告警, 1-输出告警
	LinkAlarm int `json:"libLinkAlarm"`
}
