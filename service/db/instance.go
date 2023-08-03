package db

import "sync"

var instances = sync.Map{}
var clusterInstances = sync.Map{}

func GetInstance(key string) *DBConn {
	instance, ok := instances.Load(key)
	if !ok {
		panic("not found")
	}

	ins, ok := instance.(*DBConn)
	if !ok {
		panic("not support")
	}

	return ins
}

func GetClusterInstance(key string) *ClusterConn {
	instance, ok := clusterInstances.Load(key)
	if !ok {
		panic("not found")
	}

	ins, ok := instance.(*ClusterConn)
	if !ok {
		panic("not support")
	}

	return ins
}
