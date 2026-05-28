package model

var farmDatabaseFallback = map[string]string{
	"FY": "fuyu",
	"YS": "yushu",
}

var DeviceTypeStableFallback = map[string]string{
	"STM":  "strain",
	"ACC":  "accel",
	"ACCX": "accel",
	"ACCY": "accel",
	"INS":  "inclinometer",
	"INSX": "inclinometer",
	"INSY": "inclinometer",
	"ATS":  "tension",
	"WPR":  "radar",
	"JMT":  "joint_meter",
	"HLS":  "hydrostatic",
	"ULS":  "ultrasonic_level",
	"GNSS": "gnss",
}

var stableFieldsFallback = map[string][]string{
	"strain":           {"strain"},
	"accel":            {"accel"},
	"inclinometer":     {"x", "y"},
	"tension":          {"tension"},
	"radar":            {"d", "rws", "veer", "raws", "ti", "hw_shub", "direction_hub", "v_sheer", "h_sheer", "hw_shigh", "direction_high", "hw_slow", "direction_low"},
	"joint_meter":      {"joint"},
	"hydrostatic":      {"settlement", "temperature", "pressure"},
	"ultrasonic_level": {"height"},
	"gnss":             {"longitude", "latitude", "vertical", "horizontal", "ordinate", "vertical_offset", "horizontal_offset", "ordinate_offset"},
}

var alarmDeviceTypeFallback = map[string]int64{
	"STM":  1,
	"ACC":  2,
	"ACCX": 2,
	"ACCY": 2,
	"INS":  3,
	"INSX": 3,
	"INSY": 3,
	"ATS":  4,
	"WPR":  5,
	"JMT":  6,
	"HLS":  7,
	"ULS":  8,
	"GNSS": 9,
	"IPC":  10,
}
