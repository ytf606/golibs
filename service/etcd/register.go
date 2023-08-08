package etcd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type EtcdRegister struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	prefix        string
	localAddr     string
}

func NewEtcdRegister(etcdAddr []string, prefix, localAddr string, ttl, connectTimeout int) (*EtcdRegister, error) {
	tag := "[NewEtcdRegister]"
	conf := clientv3.Config{
		Endpoints:         etcdAddr,
		DialTimeout:       time.Duration(connectTimeout) * time.Second,
		DialKeepAliveTime: time.Duration(connectTimeout) * time.Second,
	}
	ctx := context.Background()
	etcdClient, err := clientv3.New(conf)
	if err != nil {
		logx.Ex(ctx, tag, "init etcd register clientv3 failed err:%+v, conf:%+v", err, conf)
		return nil, err
	}

	// 创建一个租约
	lease := clientv3.NewLease(etcdClient)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(connectTimeout))
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, int64(ttl))
	if err != nil {
		return nil, err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return nil, err
	}
	res := &EtcdRegister{
		client:        etcdClient,
		prefix:        prefix,
		localAddr:     localAddr,
		lease:         lease,
		leaseResp:     leaseResp,
		cancelFunc:    cancel,
		keepAliveChan: leaseChannel,
	}
	go res.listenLeaseRespChan(ctx)
	return res, nil
}

func (e *EtcdRegister) listenLeaseRespChan(ctx context.Context) {
	tag := "[listenLeaseRespChan]"
	for {
		select {
		case leaseKeepResp := <-e.keepAliveChan:
			if leaseKeepResp == nil {
				logx.Ex(ctx, tag, "keep alive lease failed")
				return
			} else {
				logx.Ix(ctx, tag, "lease succ lease ID:%+v", strconv.FormatInt(int64(e.leaseResp.ID), 16))
			}
		}
	}
}

func (e *EtcdRegister) Register(ctx context.Context, service string) error {
	tag := "[Register]"
	em, err := endpoints.NewManager(e.client, e.prefix)
	if err != nil {
		return err
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	key := fmt.Sprintf("%s/%s/%s", e.prefix, service, e.localAddr)
	if err := em.AddEndpoint(ctx2, key, endpoints.Endpoint{
		Addr: e.localAddr,
	}, clientv3.WithLease(e.leaseResp.ID)); err != nil {
		logx.Ex(ctx, tag, "register etcd service failed key:%s, err:%+v", key, err)
		return err
	}
	logx.Dx(ctx, tag, "register etcd succ key:%s", key)
	return nil
}
