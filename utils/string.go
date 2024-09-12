package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RandomAlphaDigit = 0
	RandomUpperAlpha = 1
	RandomLowerAlpha = 2
	RandomDigit      = 3
)

func CreateSmsCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)) //这里面前面的04v是和后面的1000相对应的
}

func RandNumberCode() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000) //这里面前面的04v是和后面的1000相对应的
}

func DiffString(text1, text2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	return diffs
}

func IsNumber(str string) bool {
	isNumber := regexp.MustCompile(`^[0-9]*$`)
	return isNumber.MatchString(str)
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// 生成指定长度的随机字符串
func GetRandomStr(length, mode int) (string, error) {
	if length < 0 {
		return "", errors.New("illegal argument length")
	}

	upperAlphaList := "MNVWXYZABCDEFGHIJKLOPQRSTU"
	lowerAlphaList := "mnvwxyzabcdefghijklopqrstu"
	digitList := "0123456789"

	var charList string
	switch mode {
	case RandomAlphaDigit:
		charList = upperAlphaList + lowerAlphaList + digitList
	case RandomUpperAlpha:
		charList = upperAlphaList
	case RandomLowerAlpha:
		charList = lowerAlphaList
	case RandomDigit:
		charList = digitList
	default:
		return "", errors.New("illegal argument mode")
	}

	var res string
	for i := 0; i < length; i++ {
		// 生成一定范围的随机数
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(charList))
		res += string(charList[index])
	}

	return res, nil
}

// 子串
func Substring(source string, start int, end int) string {
	var des = []rune(source)
	length := len(des)

	if start < 0 || end > length || start > end {
		return ""
	}

	return string(des[start:end])
}

// 字符串长度
func LengthString(s string) int {
	return len([]rune(s))
}

// 是否为空字符串
func EmptyString(s string) bool {
	return s == ""
}

// 字符串是否相同
func EqualString(a, b string) bool {
	return a == b
}

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func NumberArrayToJoin(numbers []int64, spilt string) string {

	var temp = make([]string, len(numbers))
	for k, v := range numbers {
		temp[k] = fmt.Sprintf("%d", v)
	}
	return strings.Join(temp, spilt)
}

func StringToInt64Array(arrayString, split string) (res []int64) {
	for _, numString := range strings.Split(arrayString, split) {
		if num, err := strconv.ParseInt(numString, 10, 64); err == nil {
			res = append(res, num)
		}
	}
	return res
}

// similar_text()
func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

// levenshtein()
// costIns: Defines the cost of insertion.
// costRep: Defines the cost of replacement.
// costDel: Defines the cost of deletion.
func Levenshtein(str1, str2 string, costIns, costRep, costDel int) int {
	var maxLen = 255
	l1 := len(str1)
	l2 := len(str2)
	if l1 == 0 {
		return l2 * costIns
	}
	if l2 == 0 {
		return l1 * costDel
	}
	if l1 > maxLen || l2 > maxLen {
		return -1
	}

	tmp := make([]int, l2+1)
	p1 := make([]int, l2+1)
	p2 := make([]int, l2+1)
	var c0, c1, c2 int
	var i1, i2 int
	for i2 := 0; i2 <= l2; i2++ {
		p1[i2] = i2 * costIns
	}
	for i1 = 0; i1 < l1; i1++ {
		p2[0] = p1[0] + costDel
		for i2 = 0; i2 < l2; i2++ {
			if str1[i1] == str2[i2] {
				c0 = p1[i2]
			} else {
				c0 = p1[i2] + costRep
			}
			c1 = p1[i2+1] + costDel
			if c1 < c0 {
				c0 = c1
			}
			c2 = p2[i2] + costIns
			if c2 < c0 {
				c0 = c2
			}
			p2[i2+1] = c0
		}
		tmp = p1
		p1 = p2
		p2 = tmp
	}
	c0 = p1[l2]

	return c0
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GeneratePassword(password string) (string, string) {
	salt := GenValidateCode(6)
	password = password + salt
	return Md5ToString(password), salt
}

func CheckPassword(password, salt, checkPassword string) bool {
	password = password + salt
	check := Md5ToString(password)
	if check == checkPassword {
		return true
	}
	return false
}

func Md5ToString(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MakeToken() string {
	if u1, err := uuid.NewV4(); err != nil {
		return ""
	} else {
		return Md5ToString(u1.String())
	}
}

func ShowSubstr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func InterfaceToInt(i interface{}) (no int) {
	switch i.(type) {
	case string:
		no, _ = strconv.Atoi(i.(string))
	case int:
		no = i.(int)
		break
	}

	return no
}

func Config2Maps(str string) map[string]string {
	maps := make(map[string]string)
	spaceRe, _ := regexp.Compile("[,;\\r\\n]+")
	res := spaceRe.Split(strings.Trim(str, ",;\r\n"), -1)
	if strings.Index(str, ":") > 0 {
		for _, v := range res {
			if countSplit := strings.Split(v, ":"); len(countSplit) == 2 {
				maps[countSplit[0]] = countSplit[1]
			}
		}
	}
	return maps
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
