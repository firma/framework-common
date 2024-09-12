package utils

import (
	"fmt"
	"testing"
	"time"
)

// GetDistance 返回单位为：千米
func TestGetDistance(t *testing.T) {
	//{"plan_id":1,"longitude":120.14595106336806,"latitude":30.306600206163193}
	lat := 30.3237390
	long := 120.1215880

	lat1 := 30.306590440538194
	long1 := 120.14611382378472

	info := GetDistance(lat, long, lat1, long1)
	fmt.Println(info)
}

//func TestUint8(t *testing.T) {
//	str := "001135512"
//	data := uint8(512)
//	fmt.Println([]byte(str))
//}

func TestNames(t *testing.T) {
	timeStr := "2023-10-31"
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	fmt.Println(t2.Unix(), t2.Format("2006-01-02 15:04:05"))
}
