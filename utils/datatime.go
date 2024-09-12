package utils

import (
	"strconv"
	"time"
)

func StartOfDay(at time.Time) time.Time {
	return time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, at.Location())
}

func GetTodayString() string {
	timeStr := time.Now().Format("20060102")
	return timeStr
}

func GetYesterdayString() string {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	yesterday := t2.AddDate(0, 0, -1)

	return yesterday.Format("20060102")
}

func GetAddDay(day int) time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t2.AddDate(0, 0, day)
}

func GetTomorrowTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t2.AddDate(0, 0, 1)
}

func GetNextMonthFirstDay() time.Time {
	timeStr := time.Now().Format("2006-01")
	t2, _ := time.ParseInLocation("2006-01", timeStr, time.Local)
	return t2.AddDate(0, 1, 0)
}

func GetCurrentMonth() int64 {
	timeStr := time.Now().Format("200601")

	if month, err := strconv.Atoi(timeStr); err == nil {
		return int64(month)
	}
	return 0
}
func GetCurrentDay() int64 {
	timeStr := time.Now().Format("20060102")
	if month, err := strconv.Atoi(timeStr); err == nil {
		return int64(month)
	}
	return 0
}

func GetCurrentDaySubDayTime(subDay int) int64 {
	currentTime := time.Now()

	timeStr := currentTime.AddDate(0, 0, -subDay).Format("20060102")
	if month, err := strconv.Atoi(timeStr); err == nil {
		return int64(month)
	}
	return 0
}
func GetMonthFirstDay() time.Time {
	timeStr := time.Now().Format("2006-01")
	t2, _ := time.ParseInLocation("2006-01", timeStr, time.Local)
	return t2
}

func GetMonthStringTime(month string) time.Time {
	t2, _ := time.ParseInLocation("200601", month, time.Local)
	return t2
}

func GetTomorrowTimeString() string {
	return GetTomorrowTime().Format("2006-01-02 15:04:05")
}

func GetYesterdayTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t2.AddDate(0, 0, -1)
}

// @params today 获取日期 2023-12-12
func GetTodayStartAndEnd(day string) (*time.Time, *time.Time, error) {
	t2, err := time.ParseInLocation("2006-01-02", day, time.Local)
	if err != nil {
		return nil, nil, err
	}
	start := t2.AddDate(0, 0, 0)
	end := t2.AddDate(0, 0, 1)
	return &start, &end, err
}

func GetTodayStartTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t2.AddDate(0, 0, 0)
}

func GetTodayEndTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t2.AddDate(0, 0, 1)
}

func GetYesterdayTimeString() string {
	return GetYesterdayTime().Format("2006-01-02 15:04:05")
}

func GetYesterdayTimeKey(key string) string {
	return GetYesterdayTime().Format("2006-01-02") + ":" + key
}

func GetTodayTimeKey(key string) string {
	return time.Now().Format("2006-01-02") + ":" + key
}

func FormatWechatPay(timeStr string) time.Time {
	t2, err := time.ParseInLocation("20060102150405", timeStr, time.Local)
	if err != nil {
		return t2
	}
	return t2
}

func TimeStringSubNowDay(timeStr string) int {
	t2, _ := time.ParseInLocation("2006-01-02 15-04-05", timeStr, time.Local)
	subTime := time.Now()
	return HoursToDay(subTime, t2)
}

func TimeNowSubDay(t2 time.Time) int {
	subTime := time.Now()
	return HoursToDay(t2, subTime)
}

func HoursToDay(t1, t2 time.Time) int {
	if t1.Location().String() != t2.Location().String() {
		return 0
	}
	hours := t1.Sub(t2).Hours()
	if hours <= 0 {
		return 0
	}
	if hours < 24 {
		// may same day
		t1y, t1m, t1d := t1.Date()
		t2y, t2m, t2d := t2.Date()
		isSameDay := (t1y == t2y && t1m == t2m && t1d == t2d)
		if isSameDay {
			return 0
		} else {
			return 1
		}
	} else { // equal or more than 24
		if (hours/24)-float64(int(hours/24)) == 0 { // just 24's times
			return int(hours / 24)
		} else { // more than 24 hours
			return int(hours/24) + 1
		}
	}
}

func TimeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
