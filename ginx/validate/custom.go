package validate

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type customValidate struct {
	fn         func(validator.FieldLevel) bool
	translator string
}

//定义映射关系
var customMap = map[string]customValidate{
	"after_date": {
		fn:         afterDate,
		translator: "{0}必须要晚于当前日期",
	},
	"before_date": {
		fn:         beforeDate,
		translator: "{0}必须要早于当前日期",
	},
	"check_mobile": {
		fn:         checkMobile,
		translator: "{0}必须是一个有效的手机号码",
	},
	"check_login": {
		fn:         checkLogin,
		translator: "{0}必须由数字、字母或特殊字符组成",
	},
}

func afterDate(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	if date.Before(time.Now()) {
		return false
	}
	return true
}

func beforeDate(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	if date.After(time.Now()) {
		return false
	}
	return true
}

func checkMobile(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`^(1[3-9][0-9]\d{8})$`, fl.Field().String())
	return match
}

func checkLogin(fl validator.FieldLevel) bool {
	pwdPattern := `^[0-9a-zA-Z!@#$%^&*-_+~]+$`
	reg, err := regexp.Compile(pwdPattern) // filter exclude chars
	if err != nil {
		return false
	}

	match := reg.MatchString(fl.Field().String())
	if !match {
		return false
	}
	return true
}
