package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

func H16To10(hexStr string) int64 {
	hexNum, _ := strconv.ParseInt(hexStr, 16, 64)
	return hexNum
}

func HexToString(hexString string) string {
	data, _ := hex.DecodeString(hexString)
	result := int(data[0]) | int(data[1])<<8 | int(data[2])<<16 | int(data[3])<<24 | int(data[4])<<32 | int(data[5])<<40
	return fmt.Sprintf("%d", result)
}

func DeviceNo(data []byte) uint16 {
	//hexStr1 := "0A"
	//fmt.Println(data)
	//hexStr2 := "00"
	hexStr1 := HexTo16String(data[0:1])
	hexStr2 := HexTo16String(data[1:])

	// 将十六进制字符串解析为整数
	num1, _ := strconv.ParseUint(hexStr1, 16, 8)
	num2, _ := strconv.ParseUint(hexStr2, 16, 8)

	// 按照低位在前的顺序合并为一个 16 位整数
	result := uint16(num1) | (uint16(num2) << 8)

	//fmt.Printf("合并后的结果为 %d\n", result)
	return result
}
func HexToChart(hexStr string) string {
	decoded, _ := hex.DecodeString(hexStr)
	fmt.Println(string(decoded))
	return string(decoded)
}
func HexTo16String(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data[:len(data)] {
		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

func Hex2Dec(val string) int64 {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int64(n)
}
func HexToBytes(hexString string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Hex16to2(hex string) string {
	i, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		panic(err)
	}
	binary := fmt.Sprintf("%08b", i)
	//fmt.Printf("%s的二进制表示：%s", hex, binary)
	return binary
}

func NumberStringToArray(str string) []int64 {
	arr := make([]int64, len(str))

	for i := 0; i < len(str); i++ {
		arr[i], _ = strconv.ParseInt(string(str[i]), 10, 64) // 将字符转换为字符串
	}
	return arr
}
