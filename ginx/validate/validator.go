package validate

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

type TransLang string

const (
	ZH TransLang = "zh"
	EN TransLang = "en"
)

func InitTrans(lang TransLang) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
			if name == "" {
				name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})

		HookRegisterValidation(v)

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		trans, _ = uni.GetTranslator(string(lang))
		switch lang {
		case ZH:
			err = zh_translations.RegisterDefaultTranslations(v, trans)
			break
		case EN:
			err = en_translations.RegisterDefaultTranslations(v, trans)
			break
		default:
			err = en_translations.RegisterDefaultTranslations(v, trans)
			break
		}
		if err != nil {
			return
		}
		err = HookRegisterTranslator(v)
	}
	return
}

func GetTrans() ut.Translator {
	return trans
}
