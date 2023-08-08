package etcd

import (
	"context"
	"fmt"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	gresolver "google.golang.org/grpc/resolver"
)

type EtcdDiscover struct {
	client         *clientv3.Client
	resolverClient gresolver.Builder
	prefix         string
	connectTimeout int
	policy         string
}

func NewEtcdDiscover(etcdAddr []string, prefix string, connectTimeout int) (*EtcdDiscover, error) {
	tag := "[NewEtcdDiscover]"
	conf := clientv3.Config{
		Endpoints:         etcdAddr,
		DialTimeout:       time.Duration(connectTimeout) * time.Second,
		DialKeepAliveTime: time.Duration(connectTimeout) * time.Second,
	}
	ctx := context.Background()
	etcdClient, err := clientv3.New(conf)
	if err != nil {
		logx.Ex(ctx, tag, "init etcd discover clientv3 failed err:%+v, conf:%+v", err, conf)
		return nil, err
	}
	resolverClient, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		logx.Ex(ctx, tag, "init resolver client failed err:%+v")
		return nil, err
	}
	res := &EtcdDiscover{
		client:         etcdClient,
		resolverClient: resolverClient,
		prefix:         prefix,
		connectTimeout: connectTimeout,
		policy:         `{"loadBalancingPolicy":"round_robin"}`,
	}
	return res, nil
}

func (e *EtcdDiscover) GetService(service string) (*grpc.ClientConn, error) {
	tag := "[GetService]"
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(e.connectTimeout)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, e.GetKey(service),
		grpc.WithResolvers(e.resolverClient),
		grpc.WithDefaultServiceConfig(e.policy),
		grpc.WithInsecure(),
	)

	if err != nil {
		logx.Ex(ctx, tag, "dial remote service failed err:%+v", err)
		return nil, err
	}
	return conn, nil
}

func (e *EtcdDiscover) GetKey(service string) string {
	return fmt.Sprintf("etcd:///%s/%s", e.prefix, service)
}
