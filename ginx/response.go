package ginx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ytf606/golibs/errorx"
	"github.com/ytf606/golibs/ginx/validate"
	"github.com/ytf606/golibs/logx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func ToResponse(r *errorx.Response) *Response {
	res := &Response{
		Code: r.Code,
		Msg:  r.Message,
		Data: r.Data,
	}
	if isMsgKey == false {
		res.Message = r.Message
		res.Msg = ""
	} else {
		res.Msg = r.Message
	}
	return res
}

func ParseJson(c *gin.Context, obj interface{}, code int) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return ErrValidate(err, code)
	}
	return nil
}

func ParseCheckJson(c *gin.Context, obj interface{}, code int) error {
	if err := ParseJson(c, obj, code); err != nil {
		ErrResponse(c, err)
		return err
	}
	return nil
}

func ParseQuery(c *gin.Context, obj interface{}, code int) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return ErrValidate(err, code)
	}
	return nil
}

func ParseCheckQuery(c *gin.Context, obj interface{}, code int) error {
	if err := ParseQuery(c, obj, code); err != nil {
		ErrResponse(c, err)
		return err
	}
	return nil
}

func ParseForm(c *gin.Context, obj interface{}, code int) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return ErrValidate(err, code)
	}
	return nil
}

func ParseCheckForm(c *gin.Context, obj interface{}, code int) error {
	if err := ParseForm(c, obj, code); err != nil {
		ErrResponse(c, err)
		return err
	}
	return nil
}

func ErrValidate(err error, code int) error {
	t, ok := err.(validate.ValidationErrors)
	if ok {
		var errs validate.ValidErrors
		for key, value := range t.Translate(validate.GetTrans()) {
			errs = append(errs, &validate.ValidError{
				Key: key,
				Msg: value,
			})
		}
		err = errs
	}
	return errorx.Wrap400Response(err, code, fmt.Sprintf("parse param error: %s", err.Error()))
}

func ErrResponse(c *gin.Context, err error, status ...int) {
	var res *errorx.Response
	if res = errorx.UnWrapResponse(err); res == nil {
		res = errorx.UnWrapResponse(
			errorx.New500Response(errorx.GinxResponseTypeErr, "gin response error raw err:%+v", err),
		)
		if err != nil {
			res.Err = err
		}
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	} else if defaultHttpStatus > 0 {
		res.StatusCode = defaultHttpStatus
	}

	if err := res.Err; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}
	}

	ResJson(c, res.StatusCode, res)
}

func SuccResponse(c *gin.Context, v interface{}) {
	ResJson(c, http.StatusOK, v)
}

func ResJson(c *gin.Context, status int, v interface{}) {
	tag := "[ginx_response]"
	ctx := StdCtx(c)

	o, ok := v.(*errorx.Response)
	if ok {
		i, _ := json.Marshal(o)
		logx.Ex(ctx, tag, "response error data:%+v", string(i))
	} else {
		o = errorx.SuccResponse(v)
		i, _ := json.Marshal(o)
		logx.Dx(ctx, tag, "response success data:%+v", string(i))
	}
	c.JSON(status, ToResponse(o))
}
