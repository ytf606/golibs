package redis

import (
	"context"
	"testing"
)

func Test_Manager(t *testing.T) {
	err := DefaultManager.Init(Config{
		Name:     "test",
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
		DB:       0,
		PoolSize: 5,
	})

	if err != nil {
		t.Fatal(err)
	}

	c := DefaultManager.Get("test")

	res, err := c.Set(context.Background(), "test-key", 1, 0).Result()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	res, err = c.Get(context.Background(), "test-key").Result()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
