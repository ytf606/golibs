package request

import (
	"context"
	"strings"
	"sync"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/errorx"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx/logtrace"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
)

var (
	EtcdClientInstance sync.Map
)

type RpcxConfig struct {
	Addr           []string
	BasePath       string
	Retries        int
	ConnectTimeout time.Duration
	SerializeType  protocol.SerializeType
}

func NewRpcx(addr []string, basePath string) *RpcxConfig {
	return &RpcxConfig{
		Addr:           addr,
		BasePath:       basePath,
		Retries:        3,
		ConnectTimeout: 3 * time.Second,
		SerializeType:  protocol.JSON,
	}
}

func (c *RpcxConfig) InitEtcdClient(ctx context.Context, servicePath string) (xclient client.XClient, err error) {
	val, ok := EtcdClientInstance.Load(servicePath)
	if !ok {
		opt := client.DefaultOption
		opt.Retries = c.Retries
		opt.ConnectTimeout = c.ConnectTimeout
		opt.SerializeType = c.SerializeType
		d, err := etcdClient.NewEtcdV3Discovery(c.BasePath, servicePath, c.Addr, false, nil)
		if err != nil {
			return nil, errorx.Wrapf(err, "init etcd client failed basePath:%s, servicePath:%s, addr:%+v",
				c.BasePath, servicePath, c.Addr)
		}
		xc := client.NewXClient(servicePath, client.Failover, client.RandomSelect, d, opt)
		EtcdClientInstance.Store(servicePath, xc)
		xclient = xc
	} else {
		xclient = val.(client.XClient)
	}
	return xclient, nil
}

func (c *RpcxConfig) RpcxRequest(ctx context.Context, serviceName string, serviceMethod string, args interface{}) (*errorx.Response, error) {
	tag := "[request_rpcx_RpcxRequest]"
	reply := &errorx.Response{}
	ctx = c.GenMetadata(ctx)
	xc, err := c.InitEtcdClient(ctx, serviceName)
	if err != nil {
		logx.Ex(ctx, tag, "rpcx init etcd client failed err:%+v, serviceName:%s", err, serviceName)
		return nil, errorx.Wrap500Response(err, errorx.RpcEndpointErr, "")
	}
	tmCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	call, err := xc.Go(tmCtx, serviceMethod, args, reply, nil)
	if err != nil {
		logx.Ex(ctx, tag, "rpcx request selectClient failed err:%+v, serviceName:%s, serviceMethod:%s, args:%+v",
			err, serviceName, serviceMethod, args)
		return nil, errorx.Wrap500Response(err, errorx.RpcEndpointErr, "")
	}

	replyCall := <-call.Done
	if replyCall.Error != nil {
		logx.Ex(ctx, tag, "rpcx request replyCall failed err:%+v", replyCall)
		if err := errorx.UnWrapResponse(replyCall.Error); err != nil {
			return reply, err
		}
		return reply, errorx.Wrap500Response(replyCall.Error, errorx.RpcReturnErr, "")
	}

	if reply.Code != errorx.SuccessCode {
		//TODO 临时兼容viper中rpcx返回为0的情况
		if reply.Code == 0 && reply.Err == nil {
			return reply, nil
		}
		logx.Ex(ctx, tag, "rpcx request replyCode failed serviceName:%s, serviceMethod:%s, args:%+v, reply:%+v",
			serviceName, serviceMethod, args, reply)
		if err := errorx.UnWrapResponse(reply); err != nil {
			return reply, err
		}
		return reply,
			errorx.New500Response(errorx.RpcCodeErr,
				"rpcx request reply struct error serviceName:%s, serviceMethod:%s, args:%+v, reply:%+v",
				serviceName, serviceMethod, args, reply)
	}
	return reply, nil
}

func (c *RpcxConfig) GenMetadata(ctx context.Context) context.Context {
	var xRpcId string
	var xTraceId string
	logTraceKey := ctx.Value(logtrace.GetMetadataKey())
	if logTraceKey != nil {
		xTraceId = logTraceKey.(*logtrace.TraceNode).Get("x_trace_id")
		xTraceId = strings.Trim(xTraceId, "\"")
		xRpcId = logTraceKey.(*logtrace.TraceNode).Get("x_rpcid")
		xRpcId = strings.Trim(xRpcId, "\"")
	}
	reqMetaData := map[string]string{
		"x_trace_id": xTraceId,
		"x_rpcid":    xRpcId,
	}
	return context.WithValue(ctx, share.ReqMetaDataKey, reqMetaData)
}
