package wind

import "time"

// CalcDurationSeconds 计算两个时间戳之间的秒数，解析失败时返回 0。
func CalcDurationSeconds(startTime, endTime string) float64 {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
	}
	parse := func(s string) (time.Time, bool) {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}
	st, ok1 := parse(startTime)
	et, ok2 := parse(endTime)
	if !ok1 || !ok2 {
		return 0
	}
	d := et.Sub(st).Seconds()
	if d < 0 {
		return 0
	}
	return d
}
