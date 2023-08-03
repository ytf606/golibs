package redis

// redis manager
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	DefaultManager *RedisManager
	once           sync.Once
	Nil            = redis.Nil
)

type (
	Ins      = redis.Client
	Z        = redis.Z
	ZRangeBy = redis.ZRangeBy
)

type Config struct {
	Name           string // instance name
	Addr           string // host:port address.
	Username       string // username
	Password       string // password
	DB             int    // selected db
	PoolSize       int    // connection pool size, must > 3
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
}

type RedisManager struct {
	cfgs    map[string]Config
	servers map[string]*redis.Client
}

func NewRedisManager() *RedisManager {
	once.Do(func() {
		DefaultManager = &RedisManager{
			cfgs:    make(map[string]Config),
			servers: make(map[string]*redis.Client),
		}
	})
	return DefaultManager
}

func (r *RedisManager) Get(name string) *redis.Client {
	s, ok := r.servers[name]
	if !ok {
		panic(fmt.Sprintf("server name [%s] not found", name))
	}
	return s
}

// Init init redis manager
func (r *RedisManager) Init(configs ...Config) error {
	for _, v := range configs {
		if v.PoolSize < 3 {
			v.PoolSize = 3
		}
		if _, ok := r.cfgs[v.Name]; ok {
			return fmt.Errorf("duplicate server name:%s", v.Name)
		}
		r.cfgs[v.Name] = v
		// init redis client
		newc := redis.NewClient(&redis.Options{
			Addr:         v.Addr,
			Username:     v.Username,
			Password:     v.Password,
			DB:           v.DB,
			DialTimeout:  time.Duration(v.ConnectTimeout) * time.Second,
			ReadTimeout:  time.Duration(v.ReadTimeout) * time.Second, // default read and write timeout is 2s
			WriteTimeout: time.Duration(v.WriteTimeout) * time.Second,
			PoolSize:     v.PoolSize,
			MinIdleConns: 3,
		})
		res, err := newc.Ping(context.Background()).Result()
		if err != nil || res != "PONG" {
			return fmt.Errorf("ping redis [%s] failed, error:%s", v.Addr, err.Error())
		}
		r.servers[v.Name] = newc
	}
	return nil
}
