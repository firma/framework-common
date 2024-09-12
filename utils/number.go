package utils

import (
	"fmt"
	"sort"
	"strconv"
)

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func NumberSliceContains(arr []int64, x int64) bool {
	sort.Slice(
		arr, func(i, j int) bool {
			return arr[i] < arr[j]
		},
	)

	index := sort.Search(
		len(arr), func(i int) bool {
			return arr[i] >= x
		},
	)

	return index < len(arr) && arr[index] == x
}
