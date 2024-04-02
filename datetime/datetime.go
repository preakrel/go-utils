package datetime

import (
	"fmt"
	"time"
)

// ChangeDateFormat 转化日期格式
func ChangeDateFormat(date time.Time, isDiy bool, fmtStr string) string {
	if !isDiy {
		switch fmtStr {
		case "0":
			return date.Format("2006-01-02")
		case "1":
			return date.Format("2006-1-2")
		case "2":
			return date.Format("2006/01/02")
		case "3":
			return date.Format("2006/1/2")
		case "4":
			weekMap := map[string]string{
				"Sunday":    "星期日",
				"Monday":    "星期一",
				"Tuesday":   "星期二",
				"Wednesday": "星期三",
				"Thursday":  "星期四",
				"Friday":    "星期五",
				"Saturday":  "星期六",
			}
			return date.Format("2006年1月2日") + ", " + weekMap[date.Weekday().String()]
		case "5":
			return date.Format("06-01-02")
		case "6":
			return date.Format("06-1-2")
		case "7":
			return date.Format("06/01/02")
		case "8":
			return date.Format("06/1/2")
		case "9":
			return date.Format("2006")
		case "10":
			return date.Format("2006-01")
		case "11":
			return date.Format("2006/01")
		case "12":
			return date.Format("06-1")
		case "13":
			return date.Format("01-02")
		case "14":
			return date.Format("1-2")
		case "15":
			return date.Format("1/2")
		default:
			return date.Format("2006-1-2")
		}
	} else {
		return date.Format(fmtStr)
	}
}

// TimeToNum Unix时间戳转换
func TimeToNum(timeStr string) (int64, error) {
	t, err := time.ParseInLocation(`2006-01-02 15:04:05`, timeStr, time.Local)
	if err != nil {
		return -1, err
	}
	return t.Unix(), nil
}

// NumToTime Unix时间戳转换
func NumToTime(timeNum int64) string {
	t := time.Unix(timeNum, 0)
	return t.Format(`2006-01-02 15:04:05`)
}

// GetTimestamp 获取时间戳
func GetTimestamp() int64 {
	timestamp := time.Now().Unix()
	return timestamp
}

// GetNowTime 获取当前时间
func GetNowTime() string {
	n := time.Now().Format(`2006-01-02 15:04:05`)
	return n
}

// Sleep 毫秒延迟
func Sleep(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Millisecond)
}

// MicSleep 微秒延迟
func MicSleep(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Microsecond)
}

// NanoSleep 纳秒延迟
func NanoSleep(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Nanosecond)
}

// timeZone 默认时区
var timeZone *time.Location

// GetTimeZone 获取时区
func GetTimeZone() *time.Location {
	if timeZone == nil {
		timeZone, _ = time.LoadLocation("Local")
	}
	return timeZone
}

// SetTimeZone 设置时区
func SetTimeZone(zone string) *time.Location {
	//loc, err := time.LoadLocation("Asia/Shanghai")
	timeZone, _ = time.LoadLocation(zone)
	return timeZone
}

func StringToDate(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02",
		"02 Jan 2006",
		"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05Z07:00", // RFC3339 without T
		"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
		"2006-01-02 15:04:05",
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
	var d time.Time
	var err error
	cstZone := GetTimeZone()
	for _, dateType := range layouts {
		if d, err = time.ParseInLocation(dateType, s, cstZone); err == nil {
			return d, nil
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}
