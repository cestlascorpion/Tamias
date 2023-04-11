package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/cestlascorpion/Tamias/core"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/configor"
)

var (
	testCache *Redis
)

func init() {
	conf := &core.Config{}
	err := configor.Load(conf, "/tmp/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	cache, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	testCache = cache
}

func TestRedis_SetUri(t *testing.T) {
	if testCache == nil {
		return
	}

	err := testCache.SetUri(context.Background(), "key", "val", 1024)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestRedis_GetUri(t *testing.T) {
	if testCache == nil {
		return
	}

	uri, ttl, err := testCache.GetUri(context.Background(), "key")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(uri, ttl)

	uri, ttl, err = testCache.GetUri(context.Background(), "no-key")
	if err != redis.Nil {
		fmt.Println(err)
		t.FailNow()
	}
}
