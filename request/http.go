package request

import (
	"context"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/errorx"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	"github.com/go-resty/resty/v2"
)

var (
	HttpClient *resty.Client = resty.New().
			SetTimeout(15 * time.Second).
			SetRetryCount(3).
			SetCloseConnection(true)
	tag string = "[request_http_request]"
)

type Requester interface {
	SetTimeout(timeout int) Requester
	SetRetryCount(count int) Requester
	GetClientInstance() *resty.Client
	Get(ctx context.Context, url string, headers map[string]string) (res []byte, err error)
	GetQuery(ctx context.Context, url string, query map[string]string, headers map[string]string) (res []byte, err error)
	PostForm(ctx context.Context, url string, body map[string]string, headers map[string]string) (res []byte, err error)
	PostRaw(ctx context.Context, url, body string, headers map[string]string) (res []byte, err error)
}

type request struct {
	client *resty.Client
}

func New() Requester {
	client := resty.New().
		SetTimeout(time.Duration(HttpDefaultTimeout) * time.Second).
		SetRetryCount(HttpDefaultRetryCount).
		SetCloseConnection(true)
	return &request{
		client: client,
	}
}

func (r *request) SetTimeout(timeout int) Requester {
	r.client.SetTimeout(time.Duration(timeout) * time.Second)
	return r
}

func (r *request) SetRetryCount(count int) Requester {
	r.client.SetRetryCount(count)
	return r
}

func (r *request) GetClientInstance() *resty.Client {
	return r.client
}

func (r *request) Get(ctx context.Context, url string, headers map[string]string) (res []byte, err error) {
	resp, err := r.client.R().
		SetHeaders(headers).
		Get(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v",
			err, url, headers)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, response:%+v, cost:%v", url, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}

func (r *request) GetQuery(ctx context.Context, url string, query map[string]string, headers map[string]string) (res []byte, err error) {
	ins := r.client.R().
		SetHeaders(headers)
	if len(query) > 0 {
		ins = ins.SetQueryParams(query)
	}
	resp, err := ins.Get(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v",
			err, url, headers)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, query:%+v, response:%+v, cost:%v", url, headers, query, resp, resp.Time())
	}
	return resp.Body(), nil
}

func (r *request) PostForm(ctx context.Context, url string, body map[string]string, headers map[string]string) (res []byte, err error) {
	resp, err := r.client.R().
		SetHeaders(headers).
		SetFormData(body).
		Post(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v, body:%+v",
			err, url, headers, body)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, body:%+v, headers:%+v, response:%+v", url, body, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, body:%+v, headers:%+v, response:%+v, cost:%v", url, body, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}

func (r *request) PostRaw(ctx context.Context, url, body string, headers map[string]string) (res []byte, err error) {
	resp, err := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		SetBody(body).
		Post(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v, body:%s", err, url, headers, body)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, response:%+v, cost:%v", url, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}

func Get(ctx context.Context, url string, headers map[string]string) (res []byte, err error) {
	resp, err := HttpClient.R().
		SetHeaders(headers).
		Get(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v",
			err, url, headers)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, response:%+v, cost:%v", url, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}

func PostForm(ctx context.Context, url string, body interface{}, headers map[string]string) (res []byte, err error) {
	resp, err := HttpClient.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v, body:%+v",
			err, url, headers, body)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, response:%+v, cost:%v", url, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}

func PostRaw(ctx context.Context, url, body string, headers map[string]string) (res []byte, err error) {
	resp, err := HttpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		SetBody(body).
		Post(url)
	if err != nil {
		logx.Ex(ctx, tag, "http request failed err:%+v, url:%s, headers:%+v, body:%s", err, url, headers, body)
		return nil, errorx.Wrap500Response(err, errorx.HttpRequestReturnErr, "")
	}

	if resp.IsError() {
		logx.Ex(ctx, tag, "http request return error url:%s, headers:%+v, response:%+v", url, headers, resp)
	} else {
		logx.Dx(ctx, tag, "http request return url:%s, headers:%+v, response:%+v, cost:%v", url, headers, resp, resp.Time())
	}
	return resp.Body(), nil
}
