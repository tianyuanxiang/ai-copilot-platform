package wind

// SamplingPolicy 描述设备的采样策略
type SamplingPolicy struct {
	// IntervalMs 采样间隔（毫秒）
	IntervalMs int64
	// SamplesPerSecond 每秒采样次数
	SamplesPerSecond float64
	// RadarIndexCount WPR 雷达距离层数量，非雷达设备为 0
	RadarIndexCount int
	// Source 策略来源说明
	Source string
}

// ResolveSamplingPolicy 根据设备类型代码返回对应采样策略。
//
// 规则：
//   - ACC/ACCX/ACCY: 50Hz（每秒 50 条）
//   - IPC（网络摄像机）: 每 2 小时 1 条
//   - WPR（激光雷达）: 每 10 秒 1 个采样周期，10 个 index_id 距离层
//   - 其他设备: 1Hz（每秒 1 条）
//
// @param deviceTypeCode 设备类型代码，如 "ACC"、"WPR"、"IPC"
// @return SamplingPolicy 采样策略
func ResolveSamplingPolicy(deviceTypeCode string) SamplingPolicy {
	switch deviceTypeCode {
	case "ACC", "ACCX", "ACCY":
		return SamplingPolicy{
			IntervalMs:       20,
			SamplesPerSecond: 50,
			RadarIndexCount:  0,
			Source:           "acc_50hz",
		}
	case "IPC":
		return SamplingPolicy{
			IntervalMs:       7200000,
			SamplesPerSecond: 1.0 / 7200.0,
			RadarIndexCount:  0,
			Source:           "ipc_2h",
		}
	case "WPR":
		return SamplingPolicy{
			IntervalMs:       10000,
			SamplesPerSecond: 0.1,
			RadarIndexCount:  10,
			Source:           "wpr_10s_10index",
		}
	default:
		return SamplingPolicy{
			IntervalMs:       1000,
			SamplesPerSecond: 1.0,
			RadarIndexCount:  0,
			Source:           "default_1hz",
		}
	}
}

// ExpectedCount 计算给定时间范围内的期望数据条数。
//
// @param durationSeconds 时间范围（秒）
// @param deviceCount 设备数量（普通传感器时使用，雷达传 1）
// @return 期望数据条数
func (p SamplingPolicy) ExpectedCount(durationSeconds float64, deviceCount int64) int64 {
	if p.RadarIndexCount > 0 {
		// 雷达：按时间范围 / 10秒 * index 数
		cycles := durationSeconds / 10.0
		return int64(cycles * float64(p.RadarIndexCount))
	}
	// 普通传感器：时间范围 * 每秒采样数 * 设备数
	return int64(durationSeconds * p.SamplesPerSecond * float64(deviceCount))
}
