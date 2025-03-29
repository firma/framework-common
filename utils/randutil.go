package utils

import (
	crand "crypto/rand"

	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters2 = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/=_-"
const numbers = "1234567890"

var (
	src = rand.NewSource(time.Now().UnixNano())
)

// GenerateAbsoluteUniqueOrderNumber 生成绝对唯一的订单号
func GenerateAbsoluteUniqueOrderNumber(prefix string, refId int64) string {
	rand.Seed(time.Now().UnixNano())

	ref := UserIdToInviteCode(int(refId), 8)
	orderNumber := strings.ToUpper(fmt.Sprintf("%s%s%010d", ref, prefix, rand.Intn(9999999999)))
	return orderNumber
}

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandStr2(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters2) {
			b[i] = letters2[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandNumber(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(numbers) {
			b[i] = numbers[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func GenerateCaptchaCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	// 使用新的随机数源创建一个随机数生成器
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var strBuilder strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&strBuilder, "%d", numeric[rnd.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return strBuilder.String()
}

func GenerateUid() string {
	var uidlen int = 6
	var uidpadding int64 = int64(math.Pow(10, float64(uidlen-1)))
	milli := time.Now().UnixMilli()
	millistr := fmt.Sprintf("%d", milli)
	millistr = millistr[len(millistr)-uidlen:]
	if strings.HasPrefix(millistr, "0") {
		millistr = strconv.FormatInt(milli+uidpadding, 10)
	}
	millistr = millistr[len(millistr)-uidlen:]
	uid := fmt.Sprintf("%s%s", millistr, RandNumber(4))
	return uid
}

func GenerateUidNew() string {
	var uidlen int = 7
	milli := time.Now().Unix()
	millistr := fmt.Sprintf("%d", milli)
	millistr = millistr[:uidlen]

	randnumber := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	uid := fmt.Sprintf("%s%s", millistr, randnumber)
	return uid
}

func RandString(size int) string {
	key := make([]byte, size)
	//rand.Read(key)
	if _, err := io.ReadFull(crand.Reader, key); err != nil {
		return ""
	}
	return hex.EncodeToString(key)
}

func GenerateRandomCode(numDigits, numLetters int) string {
	const digits = "0123456789"
	const lowerCaseLetters = "abcdefghijkmnpqrstuvwxyz"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	totalLength := numDigits + numLetters
	for {
		b := make([]byte, totalLength)
		for i := 0; i < numDigits; i++ {
			time.Sleep(1 * time.Nanosecond)
			b[i] = digits[seededRand.Intn(len(digits))]
		}
		for i := numDigits; i < totalLength; i++ {
			time.Sleep(1 * time.Nanosecond)
			b[i] = lowerCaseLetters[seededRand.Intn(len(lowerCaseLetters))]
		}

		// 随机打乱字符顺序
		rand.Shuffle(
			totalLength, func(i, j int) {
				b[i], b[j] = b[j], b[i]
			},
		)

		return string(b)
	}
}

func GenerateRandomCodes(num, numDigits, numLetters int) []string {
	codes := make([]string, num)

	existingCodes := make(map[string]struct{})

	i := 0
	for i < num {
		randomCode := GenerateRandomCode(numDigits, numLetters)
		// 检查是否已经存在
		if _, exists := existingCodes[randomCode]; !exists {
			existingCodes[randomCode] = struct{}{}
			codes[i] = randomCode
			i++
		}
	}

	return codes
}
