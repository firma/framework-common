package validate

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/locales/zh"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	//注册翻译器
	z := zh.New()
	uni = ut.New(z, z)

	trans, _ = uni.GetTranslator("zh")

	//获取gin的校验器
	validate = binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	if err := zhtrans.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Printf("register validate translation err: %s", err)
	}
	validate.RegisterValidation("mobile", ValidateMobile)
	validate.RegisterValidation("abc", ValidateAAbc)
	validate.RegisterValidation("tell", ValidateTell)
}
