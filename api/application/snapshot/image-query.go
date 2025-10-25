/**
 * @Author: Kaicen-X
 * @Date: 2025/03/02 00:17
 * @Description: 华为HoloSens SDC API北向接口抓拍图片查询
 */
package snapshot

import (
	"net/url"
	"strconv"
)

// SnapshotType 抓拍图类型
type SnapshotType int

const (
	// 抓拍图类型：手动抓拍图
	SnapshotType_Manual SnapshotType = 0
	// 抓拍图类型：告警联动抓拍图
	SnapshotType_Alarm = 1
	// 抓拍图类型：定时抓拍图
	SnapshotType_Timer = 2
	// 抓拍图类型：目标抓拍图
	SnapshotType_FaceDT = 4
	// 抓拍图类型：机非人抓拍图
	SnapshotType_ITGT = 5
	// 抓拍图类型：ITS、违停抓拍图
	SnapshotType_ITS = 6
)

// SnapshotVehicleType 车辆类型
type SnapshotVehicleType int

// 车辆类型枚举: 可能不是完整枚举，但不影响实际值；详情请查阅官方文档
const (
	SnapshotVehicleType_NOT         SnapshotVehicleType = 0
	SnapshotVehicleType_CAR         SnapshotVehicleType = 1
	SnapshotVehicleType_TRUCK       SnapshotVehicleType = 2
	SnapshotVehicleType_VAN         SnapshotVehicleType = 3
	SnapshotVehicleType_PAS         SnapshotVehicleType = 4
	SnapshotVehicleType_BUGGY       SnapshotVehicleType = 5
	SnapshotVehicleType_SUV         SnapshotVehicleType = 6
	SnapshotVehicleType_MEDIUM_BUS  SnapshotVehicleType = 7
	SnapshotVehicleType_MOT         SnapshotVehicleType = 8
	SnapshotVehicleType_PEDESTRAIN  SnapshotVehicleType = 9
	SnapshotVehicleType_SCHOOL_BUS  SnapshotVehicleType = 10
	SnapshotVehicleType_HEAVY_TRUCK SnapshotVehicleType = 11
	SnapshotVehicleType_TANKER      SnapshotVehicleType = 12
	SnapshotVehicleType_RIDEMAN     SnapshotVehicleType = 13
	SnapshotVehicleType_CAR_M       SnapshotVehicleType = 14
	SnapshotVehicleType_CAR_L       SnapshotVehicleType = 15
	SnapshotVehicleType_CAR_S       SnapshotVehicleType = 16
	SnapshotVehicleType_CAR_TWO     SnapshotVehicleType = 17
	SnapshotVehicleType_CAR_THR     SnapshotVehicleType = 18
	SnapshotVehicleType_SUV_QINGKE  SnapshotVehicleType = 19
	SnapshotVehicleType_SUV_L       SnapshotVehicleType = 20
	SnapshotVehicleType_SUV_S       SnapshotVehicleType = 21
	SnapshotVehicleType_SUV_Z       SnapshotVehicleType = 22
	SnapshotVehicleType_SUV_M       SnapshotVehicleType = 23
	SnapshotVehicleType_SUV_B       SnapshotVehicleType = 24
	SnapshotVehicleType_WEIMIAN     SnapshotVehicleType = 25
	SnapshotVehicleType_MPV         SnapshotVehicleType = 26
	SnapshotVehicleType_JP          SnapshotVehicleType = 27
	SnapshotVehicleType_WEIKA       SnapshotVehicleType = 28
	SnapshotVehicleType_SUV_PIKA    SnapshotVehicleType = 29
	SnapshotVehicleType_TRUCK_Z     SnapshotVehicleType = 30
	SnapshotVehicleType_TRUCK_KEI   SnapshotVehicleType = 31
	SnapshotVehicleType_ZHONGKA     SnapshotVehicleType = 32
	SnapshotVehicleType_TAXI        SnapshotVehicleType = 33
	SnapshotVehicleType_TANK        SnapshotVehicleType = 34
	SnapshotVehicleType_CRANE       SnapshotVehicleType = 35
	SnapshotVehicleType_MOTOR       SnapshotVehicleType = 36
	SnapshotVehicleType_ALL         SnapshotVehicleType = 512
)

// SnapshotDeviceType 抓拍机类型
type SnapshotDeviceType int

const (
	// 抓拍机类型： 普通抓拍机
	SnapshotDeviceType_Normal SnapshotDeviceType = 0
	// 抓拍机类型： 抓拍机（微bayonet）
	SnapshotDeviceType_Microbayonet SnapshotDeviceType = 1
)

// SnapshotVehicleRegulationType 违章类型
type SnapshotVehicleRegulationType int

// 违章类型枚举: 可能不是完整枚举，但不影响实际值；详情请查阅官方文档
const (
	SnapshotVehicleRegulationType_NULL                         SnapshotVehicleRegulationType = 0
	SnapshotVehicleRegulationType_PASS_SNAP                    SnapshotVehicleRegulationType = 1
	SnapshotVehicleRegulationType_RUNNINGRED                   SnapshotVehicleRegulationType = 2
	SnapshotVehicleRegulationType_EXCEET_SPEED                 SnapshotVehicleRegulationType = 3
	SnapshotVehicleRegulationType_UNDER_SPEED                  SnapshotVehicleRegulationType = 4
	SnapshotVehicleRegulationType_WRONG_DIRECTION              SnapshotVehicleRegulationType = 5
	SnapshotVehicleRegulationType_REVERSE_DIRECTION            SnapshotVehicleRegulationType = 6
	SnapshotVehicleRegulationType_LICENCE_RESTRICTE            SnapshotVehicleRegulationType = 7
	SnapshotVehicleRegulationType_OVER_LANE_LINE               SnapshotVehicleRegulationType = 8
	SnapshotVehicleRegulationType_CHANGE_LANE                  SnapshotVehicleRegulationType = 9
	SnapshotVehicleRegulationType_MOTOR_IN_BICYCLE_LANE        SnapshotVehicleRegulationType = 10
	SnapshotVehicleRegulationType_ACCOMMODATION_LANE           SnapshotVehicleRegulationType = 11
	SnapshotVehicleRegulationType_PROHIBITION                  SnapshotVehicleRegulationType = 12
	SnapshotVehicleRegulationType_REMAIN_CROSS                 SnapshotVehicleRegulationType = 13
	SnapshotVehicleRegulationType_STOP_LIGHT_GREEN             SnapshotVehicleRegulationType = 14
	SnapshotVehicleRegulationType_EMERGENCY                    SnapshotVehicleRegulationType = 15
	SnapshotVehicleRegulationType_SAFETY_BELT                  SnapshotVehicleRegulationType = 16
	SnapshotVehicleRegulationType_U_TURN                       SnapshotVehicleRegulationType = 17
	SnapshotVehicleRegulationType_PORT_VEHICLE_DRIVER_CALL     SnapshotVehicleRegulationType = 18
	SnapshotVehicleRegulationType_PORT_VEHICLE_COPILOT_NO_BELT SnapshotVehicleRegulationType = 19
	SnapshotVehicleRegulationType_PORT_VEHICLE_NO_YEAR_LOG     SnapshotVehicleRegulationType = 20
	SnapshotVehicleRegulationType_ILLEGAL_PARKING              SnapshotVehicleRegulationType = 21
	SnapshotVehicleRegulationType_IMPOLITE_PEDESTRIANS         SnapshotVehicleRegulationType = 22
	SnapshotVehicleRegulationType_OCCUPANCY_BUSLAN             SnapshotVehicleRegulationType = 23
	SnapshotVehicleRegulationType_LARGER_VEHICLE_OUTOFLANE     SnapshotVehicleRegulationType = 24
	SnapshotVehicleRegulationType_IBALL_ILLEGAL_PARKING        SnapshotVehicleRegulationType = 25
	SnapshotVehicleRegulationType_VEHICLEBAN                   SnapshotVehicleRegulationType = 26
	SnapshotVehicleRegulationType_NOVEHICLE_IN_LANE            SnapshotVehicleRegulationType = 27
	SnapshotVehicleRegulationType_PED_RUN_RED                  SnapshotVehicleRegulationType = 28
	SnapshotVehicleRegulationType_PASS_SNAP_FAKE               SnapshotVehicleRegulationType = 29
	SnapshotVehicleRegulationType_RECOGNIZE_ONLY               SnapshotVehicleRegulationType = 30
	SnapshotVehicleRegulationType_ABNORMAL_PLATE               SnapshotVehicleRegulationType = 31
	SnapshotVehicleRegulationType_PED_SNAP                     SnapshotVehicleRegulationType = 32
	SnapshotVehicleRegulationType_NMV_SNAP                     SnapshotVehicleRegulationType = 33
	SnapshotVehicleRegulationType_MOTORCYCLE_FORBIDDEN         SnapshotVehicleRegulationType = 34
)

// QueryParam 抓拍图查询参数构建器
type QueryParam func(url.Values)

// QueryWithTotalNum 抓拍图查询参数：总记录数
//
//	通常首次不填写，第二次查询时使用第一次查询时返回的结果
func QueryWithTotalNum(val int) QueryParam {
	return func(values url.Values) {
		values.Set("TotalNum", strconv.Itoa(val))
	}
}

// QueryWithBeginIndex 抓拍图查询参数：当前开始记录数
//
//	取值范围：[1,4096]
func QueryWithBeginIndex(val int) QueryParam {
	return func(values url.Values) {
		values.Set("BeginIndex", strconv.Itoa(val))
	}
}

// QueryWithEndIndex 抓拍图查询参数：当前结束记录数
//
//	取值范围：[1,4096]
func QueryWithEndIndex(val int) QueryParam {
	return func(values url.Values) {
		values.Set("EndIndex", strconv.Itoa(val))
	}
}

// QueryWithBeginTime 抓拍图查询参数：开始时间
//
//	UTC时间戳，单位秒
func QueryWithBeginTime(val int64) QueryParam {
	return func(values url.Values) {
		values.Set("BeginTime", strconv.FormatInt(val, 10))
	}
}

// QueryWithEndTime 抓拍图查询参数：结束时间
//
//	UTC时间戳，单位秒
func QueryWithEndTime(val int64) QueryParam {
	return func(values url.Values) {
		values.Set("EndTime", strconv.FormatInt(val, 10))
	}
}

// QueryWithTimeType 抓拍图查询参数：查询时间类型
//
//	不带此字段时，默认按照UTC时间查询
//	取值范围：（0-摄像机本地时间，1-UTC时间）
func QueryWithTimeType(val int) QueryParam {
	return func(values url.Values) {
		values.Set("TimeType", strconv.Itoa(val))
	}
}

// QueryWithSnapshotType 抓拍图查询参数：抓拍图类型
func QueryWithSnapshotType(val SnapshotType) QueryParam {
	return func(values url.Values) {
		values.Set("SnapshotType", strconv.Itoa(int(val)))
	}
}

// QueryWithLaneId 抓拍图查询参数：车道号
//
//	取值范围：[0,3]
func QueryWithLaneId(val int) QueryParam {
	return func(values url.Values) {
		values.Set("LaneId", strconv.Itoa(val))
	}
}

// QueryWithVehicleType 抓拍图查询参数：车辆类型
func QueryWithVehicleType(val SnapshotVehicleType) QueryParam {
	return func(values url.Values) {
		values.Set("VehicleType", strconv.Itoa(int(val)))
	}
}

// QueryWithSnapshotDevType 抓拍图查询参数：抓拍机类型
func QueryWithSnapshotDevType(val SnapshotDeviceType) QueryParam {
	return func(values url.Values) {
		values.Set("SnapshotDevType", strconv.Itoa(int(val)))
	}
}

// QueryWithVehicleRegulationType 抓拍图查询参数：违章类型
func QueryWithVehicleRegulationType(val SnapshotVehicleRegulationType) QueryParam {
	return func(values url.Values) {
		values.Set("VehicleRegulationType", strconv.Itoa(int(val)))
	}
}

// QueryWithRegulationRecordTime 抓拍图查询参数：违章录像时间长度
//
//	取值范围：[10,20]
func QueryWithRegulationRecordTime(val int) QueryParam {
	return func(values url.Values) {
		values.Set("RegulationRecordTime", strconv.Itoa(val))
	}
}

// QueryWithOnceInquireFlag 抓拍图查询参数：单页查询标记
//
//	取值范围：（1-单页查询使能，0-单页查询关闭）
func QueryWithOnceInquireFlag(val int) QueryParam {
	return func(values url.Values) {
		values.Set("OnceInquireFlag", strconv.Itoa(val))
	}
}

// SnapshotImageInfo 抓拍图信息
type SnapshotImageInfo struct {
	// 抓拍时间（单位秒）
	SnapTime int64 `json:"snapTime"`
	// 抓拍时间类型（0-摄像机本地时间，1-UTC时间）
	TimeType int `json:"timeType"`
}

// SnapshotRecordInfo 抓拍记录信息
type SnapshotRecordInfo struct {
	// 是否存在关联录像
	RecordExist bool `json:"recordExist"`
	// 关联录像开始时间（UTC时间戳，单位秒）
	StartTime int64 `json:"startTime"`
	// 关联录像结束时间（UTC时间戳，单位秒）
	EndTime int64 `json:"endTime"`
}

// QueryReply 抓拍图查询响应
type QueryReply struct {
	// 总记录数
	TotalNum int `json:"totalNum"`
	// 当前开始记录数
	BeginIndex int `json:"beginIndex"`
	// 当前结束记录数
	EndIndex int `json:"endIndex"`
	// 抓拍图信息列表
	ImageInfoList []SnapshotImageInfo `json:"imageInfoList"`
	// 抓拍图大小（单位：字节）
	ContentSize uint32 `json:"contentSize"`
	// 抓拍图名称（长度63字符）
	ContentId string `json:"contentId"`
	// 抓拍记录信息
	RecordInfo SnapshotRecordInfo `json:"recordInfo"`
	// 车道号
	//
	//	取值范围：[0,3]
	LaneId int `json:"laneId"`
	// 车辆类型
	VehicleType SnapshotVehicleType `json:"vehicleType"`
	// 违章类型
	VehicleRegulationType SnapshotVehicleRegulationType `json:"vehicleRegulationType"`
}

// ImageQuery 抓拍图查询
//
//	@param uuid: 设备通道UUID
//	@param params: 查询参数
//	@return 查询结果
//	@return 异常信息
func (p *Manager) ImageQuery(uuid string, params ...QueryParam) (*QueryReply, error) {
	// 获取Socket连接的HTTP客户端
	client := p.connInstance.LockHttpClient()
	defer p.connInstance.Unlock()

	// 构建查询请求
	req := client.Get("/SDCAPI/V1.0/Storage/Snapshot/Inquire").
		SetContentType("application/x-www-form-urlencoded").
		AddQuery("UUID", uuid)
	// 检查参数，执行添加
	for _, param := range params {
		param(req.GetQuery())
	}
	// 发送请求
	var reply QueryReply
	_, err := req.DecodeJSON(&reply)
	if err != nil {
		return nil, err
	}

	// OK
	return &reply, nil
}
