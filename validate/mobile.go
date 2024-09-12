package validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	// 利用反射拿到结构体tag含有mobile的key字段
	mobile := fl.Field().String()
	return checkMobile(mobile)
}

func checkMobile(mobile string) bool {
	//使用正则表达式判断是否合法 /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/
	if len(mobile) != 11 {
		return false
	}
	//使用正则表达式判断是否合法 /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/
	ok, _ := regexp.MatchString(`^1(3\d|4\d|5\d|6\d|7\d|8\d|9\d|2\d|0\d|)\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
