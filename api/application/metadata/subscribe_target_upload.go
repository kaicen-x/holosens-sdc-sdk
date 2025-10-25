/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口智能元数据订阅上报
 */
package metadata

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// SubscribeUploadCommonInfo 元数据订阅目标数据上报通用信息
type SubscribeUploadCommonInfo struct {
	UUID     string `json:"UUID"`     // 通道ID
	DeviceID string `json:"deviceID"` // 设备ID
}

// SubscribeTargetUploadSubImageInfo 元数据订阅目标数据上报目标图像信息
type SubscribeTargetUploadSubImageInfo struct {
	// 图像ID，编码规则：设备编码+抓拍时间YYMMDDhhmmssSSS+目标ID（HEX编码）
	// 	其中：
	//	1、目标ID为INT64类型，经过HEX编码后进行拼接
	//	2、时间采用UTC，精确到毫秒
	ImageID string `json:"imageID"`

	// 关联的图像ID
	RelatedImageIDs []string `json:"relatedImageIDs"`

	// 图像类型
	//	取值范围：15-全景图，11-目标图，10-目标整体图
	ImageType int64 `json:"imageType"`

	// 抓拍时间（UTC）
	// 	格式：YYMMDDhhmmssSSS，
	ShotTime string `json:"shotTime"`

	// 置信度，目标抠图kps质量过滤标志位
	//	取值范围：0和1，0表示过滤抓拍，1表示正常抓拍
	PicKps int64 `json:"picKps"`

	// 图像左上角X万分比坐标
	// 	抠图对象中包含本抠图在全景图中的坐标，全景图和无坐标的抠图此字段无效
	LeftTopX int64 `json:"leftTopX"`

	// 图像左上角Y万分比坐标
	// 	抠图对象中包含本抠图在全景图中的坐标，全景图和无坐标的抠图此字段无效
	LeftTopY int64 `json:"leftTopY"`

	// 图像右下角X万分比坐标
	// 	抠图对象中包含本抠图在全景图中的坐标，全景图和无坐标的抠图此字段无效
	RightBtmX int64 `json:"rightBtmX"`

	// 图像右下角Y万分比坐标
	// 	抠图对象中包含本抠图在全景图中的坐标，全景图和无坐标的抠图此字段无效
	RightBtmY int64 `json:"rightBtmY"`

	// 图像宽度(px)
	//	目标图此字段无效
	Width int64 `json:"width"`

	// 图像高度(px)
	//	目标图此字段无效
	Height int64 `json:"height"`

	// Base64编码的抓拍图片(jpeg)
	Data string `json:"data"`
}

// SubscribeTargetUploadDetailInfo 元数据订阅目标数据上报详细信息
type SubscribeTargetUploadDetailInfo struct {
	// 元数据类型
	// 	1 - 目标抓拍
	// 	2 - 目标识别
	// 	53 - 骑行人
	TargetType int64 `json:"targetType"`

	// 目标ID，编码格式：设备编码+抓拍时间YYMMDDhhmmssSSS+目标ID（HEX编码）
	//	其中：
	//	1、目标ID为INT64类型，经过HEX编码后进行拼接
	//	2、时间采用UTC，精确到毫秒
	FaceID string `json:"faceID"`

	// 目标识别算法版本号
	// 	最长48个字符
	FaceRecAlgVersion string `json:"faceRecAlgVersion"`

	// 目标识别抠图质量分
	// 	取值范围：0~100
	QualityScore int64 `json:"qualityScore"`

	// 目标库中匹配上人员的名称
	// 	最长63个字符（字母、数字、汉字），每个汉字占3个字符
	Name string `json:"name"`

	// 目标库中匹配上的人员的证件类型
	// 	取值范围：0-身份证，1-护照，2-军官证，3-驾驶证，4-其他
	IDType int64 `json:"IDType"`

	// 目标库中匹配上的人员的证件号
	// 	最多支持31位字符（字母，数字和汉字），每个汉字占3位字符
	IDNumber string `json:"IDNumber"`

	// 目标库中匹配上的人员的出生日期
	// 	最长32位字符，每个汉字占3位字符
	Birthday string `json:"birthday"`

	// 目标库中匹配上的人员所属的省
	// 	最长32位字符，每个汉字占3位字符
	Province string `json:"province"`

	// 目标库中匹配上的人员所属的城市
	// 	最长48个字符
	City string `json:"city"`

	// 人员性别
	// 	取值范围：0-未识别，1-女，2-男
	GenderCode int64 `json:"genderCode"`

	// 年龄上限
	AgeUpLimit int64 `json:"ageUpLimit"`

	// 年龄下限
	AgeLowerLimit int64 `json:"ageLowerLimit"`

	// 发型
	// 	取值范围：0-未识别，1-长头发，2-短头发，3-秃头
	HairStyle int64 `json:"hairStyle"`

	// 遮档(口罩)
	// 	取值范围：0-未识别，1-未带口罩，2-戴口罩
	HasRespirator int64 `json:"hasRespirator"`

	// 范佩戴口罩
	// 	取值范围：0-未知，1-规范戴口罩，2-不规范戴口罩
	MouthMaskStandard int64 `json:"mouthmaskStandard"`

	// 戴帽子
	// 	取值范围：0-未识别，1-未戴帽子，2-戴帽子
	HasCap int64 `json:"hasCap"`

	// 眼镜样式
	// 	取值范围：0-未识别，1-未戴眼镜，2-戴普通眼镜，3-戴太阳眼镜
	HasGlass int64 `json:"hasGlass"`

	// 上衣款式
	// 	取值范围：0-未识别，1-长袖，2-短袖
	UpperStyle int64 `json:"upperStyle"`

	// 上衣颜色
	// 	取值范围：0-未识别，1-黑，2-蓝，3-绿，4-白/灰，5-黄/橙/棕，6-红/粉/紫
	UpperColor int64 `json:"upperColor"`

	// 上衣纹理
	// 	取值范围：0-未识别，1-纯色，2-条纹，3-格子
	UpperTexture int64 `json:"upperTexture"`

	// 下衣款式
	// 	取值范围：0-未识别，1-长裤，2-短裤，3-裙子
	LowStyle int64 `json:"lowStyle"`

	// 下衣颜色
	// 	取值范围：0-未识别，1-黑，2-蓝，3-绿，4-白/灰，5-黄/橙/棕，6-红/粉/紫
	LowerColor int64 `json:"lowerColor"`

	// 体型
	// 	取值范围：0-未识别，1-标准，2-胖，3-瘦
	BodyType int64 `json:"bodyType"`

	// 背包
	// 	取值范围：0-未识别，1-未背包，2-背包
	HasBackpack int64 `json:"hasBackpack"`

	// 胡子
	// 	取值范围：0-未识别，1-没有胡子，2-有胡子
	HasMustache int64 `json:"hasMustache"`

	// 是否拎东西
	// 	取值范围：0-未识别，1-未拎东西，2-拎东西
	CarryBag int64 `json:"carryBag"`

	// 斜挎包
	// 	取值范围：0-未识别，1-无斜挎包，2-有斜挎包
	HasSatchel int64 `json:"hasSatchel"`

	// 前面背包
	// 	取值范围：0-未识别，1-无前背包，2-有前背包
	HasFrontBag int64 `json:"hasFrontBag"`

	//
	// 	取值范围：0-未识别，1-无雨伞，2-有雨伞
	HasUmbrella int64 `json:"hasUmbrella"`

	// 行李箱
	// 	取值范围：0-未识别，1-无行李箱，2-有行李箱
	HasLuggage int64 `json:"hasLuggage"`

	// 行进方向
	// 	取值范围：0-未识别，1-朝前，2-朝后
	MoveDirection int64 `json:"moveDirection"`

	// 行进速度
	// 	取值范围：0-未识别，1-慢速，2-快速
	MoveSpeed int64 `json:"moveSpeed"`

	// 朝向
	// 	取值范围：0-未识别，1-朝前，2-朝后，3-朝左，4-朝右
	HumanView int64 `json:"humanView"`

	// 图片数据对象数组
	// 	可能包含目标抠图、目标整体抠图、全景图
	SubImageList []SubscribeTargetUploadSubImageInfo `json:"subImageList"`
}

// SubscribeTargetUploadInfo 元数据订阅目标数据上报信息
type SubscribeTargetUploadInfo struct {
	Common     SubscribeUploadCommonInfo         `json:"common"`     // 元数据通用信息
	TargetList []SubscribeTargetUploadDetailInfo `json:"targetList"` // 元数据详细信息列表
}

// SubscribeTargetUploadParams 元数据订阅目标数据上报参数
type SubscribeTargetUploadParams struct {
	Metadata SubscribeTargetUploadInfo `json:"metadataObject"` // 元数据上报信息
}

// SubscribeTargetUpload 元数据订阅目标数据上报（HTTP响应已被接管，请不要再发送任何响应）
func SubscribeTargetUpload(w http.ResponseWriter, r *http.Request) (*SubscribeTargetUploadParams, error) {
	// 检查请求体
	if r.Body == nil {
		return nil, errors.New("request body is empty")
	}
	// 延迟关闭请求体
	defer r.Body.Close()
	// 读取请求
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// 解析请求参数
	var params SubscribeTargetUploadParams
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, err
	}

	// 响应成功状态码
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	// OK
	return &params, nil
}
