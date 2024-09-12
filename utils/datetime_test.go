package utils

import (
	"fmt"
	"testing"
)

func TestBetweenTodayForStartEnd(t *testing.T) {
	start, end, err := GetTodayStartAndEnd("2023-12-12")
	fmt.Println(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), err)
}
