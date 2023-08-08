package validate

import (
	"github.com/ytf606/golibs/logx"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func HookRegisterValidation(obj *validator.Validate) error {
	// 在校验器注册自定义的校验方法
	var err error
	for k, v := range customMap {
		if err = obj.RegisterValidation(k, v.fn); err != nil {
			logx.E("[ginx_validate_HookRegisterValidation]", "register custom validate failed err:%+v, func:%s", err, k)
		}
	}
	return err
}

func HookRegisterTranslator(obj *validator.Validate) error {
	var err error
	for k, v := range customMap {
		if err = obj.RegisterTranslation(
			k,
			GetTrans(),
			registerTranslator(k, v.translator),
			translate,
		); err != nil {
			logx.E("[ginx_validate_HookRegisterTranslator]", "register custom translator failed err:%+v, func:%s", err, k)
		}
	}
	return nil
}

// registerTranslator 为自定义字段添加翻译功能
func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		logx.E("[ginx_validate_translate]", "custom translate failed err:%+v, tag:%s", err, fe.Tag())
	}
	return msg
}
