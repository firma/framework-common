package utils

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	//001914808
	var hexStr int64
	hexStr = 1914808
	fmt.Println(hexStr)
	//num := "001914808"
	//hexStr, err := strconv.ParseInt(num, 10, 64)
	//fmt.Println(hexStr)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	hex := strconv.FormatInt(hexStr, 16)

	bytes, _ := HexToBytes(hex)
	fmt.Println()

	for _, item := range bytes {

		s := strconv.FormatInt(int64(item&0xff), 16)
		i, _ := strconv.ParseInt(s, 16, 64)

		fmt.Println(HexTo16String(bytes), item, s, i, Hex2Dec(s))
	}

	snNo := 1914808

	numStr := strconv.Itoa(snNo)
	lenInfo := len(numStr)
	start := lenInfo - 3
	areaCode, _ := strconv.Atoi(numStr[start:])
	start = lenInfo - 6
	phoneNum, _ := strconv.Atoi(numStr[start:4])
	end := 3
	if lenInfo == 9 {
		end = 9 - lenInfo
	}

	lastThree := numStr[len(numStr)-3:]

	lastThree2 := numStr[len(numStr)-6:]
	fmt.Println("lastThree2", lastThree2)
	last := strings.Replace(numStr, lastThree2, "", -1)
	fmt.Println("last", numStr, lastThree2, last)
	lastThree2 = lastThree2[:3]
	//lastThree3 := numStr[len(numStr)-9:]

	phoneNums, _ := strconv.Atoi(numStr[0:end])
	fmt.Println(start, end, lastThree, lastThree2, "|", last, "|", areaCode, phoneNum, phoneNums)

}
