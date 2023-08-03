package errorx

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	SuccessStatusCode = 200
	SuccessCode       = 200
	SuccessMsg        = "ok"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Err        error       `json:"err"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (r *Response) Error() string {
	if r.Message != "" && r.Err != nil {
		return fmt.Sprintf("%s: %s", r.Message, r.Err.Error())
	}
	if r.Err != nil {
		return r.Err.Error()
	}
	return r.Message
}

func UnWrapResponse(err error) *Response {
	if v, ok := err.(*Response); ok {
		return v
	}
	return nil
}

func SuccResponse(data interface{}) *Response {
	return &Response{
		StatusCode: SuccessStatusCode,
		Code:       SuccessCode,
		Message:    SuccessMsg,
		Data:       data,
	}
}

func WrapErrResponse(err error, statusCode, code int, message string, args ...interface{}) error {
	res := &Response{
		Err:        err,
		StatusCode: statusCode,
		Code:       code,
		Message:    fmt.Sprintf(message, args...),
	}
	return res
}

func NewErrResponse(statusCode, code int, message string, args ...interface{}) error {
	res := &Response{
		Code:       code,
		Message:    fmt.Sprintf(message, args...),
		StatusCode: statusCode,
	}
	res.Err = errors.New(res.Message)
	return res
}

func New400Response(code int, msg string, args ...interface{}) error {
	return NewErrResponse(400, code, msg, args...)
}

func New401Response(code int, msg string, args ...interface{}) error {
	return NewErrResponse(401, code, msg, args...)
}

func New403Response(code int, msg string, args ...interface{}) error {
	return NewErrResponse(403, code, msg, args...)
}

func New404Response(code int, msg string, args ...interface{}) error {
	return NewErrResponse(404, code, msg, args...)
}

func New500Response(code int, msg string, args ...interface{}) error {
	return NewErrResponse(500, code, msg, args...)
}

func Wrap400Response(err error, code int, msg string, args ...interface{}) error {
	return WrapErrResponse(err, 400, code, msg, args...)
}

func Wrap401Response(err error, code int, msg string, args ...interface{}) error {
	return WrapErrResponse(err, 401, code, msg, args...)
}

func Wrap403Response(err error, code int, msg string, args ...interface{}) error {
	return WrapErrResponse(err, 403, code, msg, args...)
}

func Wrap404Response(err error, code int, msg string, args ...interface{}) error {
	return WrapErrResponse(err, 404, code, msg, args...)
}

func Wrap500Response(err error, code int, msg string, args ...interface{}) error {
	return WrapErrResponse(err, 500, code, msg, args...)
}
