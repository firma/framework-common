package validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 校验手机号
func ValidateAAbc(fl validator.FieldLevel) bool {
	// 利用反射拿到结构体tag含有mobile的key字段

	abc := fl.Field().String()
	return checkAbc(abc)
}

func checkAbc(abc string) bool {
	//使用正则表达式判断是否合法 /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/
	ok, _ := regexp.MatchString(`^[A-Za-z]+$`, abc)
	if !ok {
		return false
	}
	return true
}

// ValidateTell 校验电话
func ValidateTell(fl validator.FieldLevel) bool {
	// 利用反射拿到结构体tag含有mobile的key字段
	tell := fl.Field().String()
	if checkMobile(tell) == false {
		return checkTell(tell)
	} else {
		return true
	}
}

func checkTell(abc string) bool {
	//使用正则表达式判断是否合法 /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/
	ok, _ := regexp.MatchString(`^[-|0-9]+$`, abc)
	if !ok {
		return false
	}
	return true
}

func checkNumber(abc string) bool {
	//使用正则表达式判断是否合法 /^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$/
	ok, _ := regexp.MatchString(`^[-|0-9]+$`, abc)
	if !ok {
		return false
	}
	return true
}
