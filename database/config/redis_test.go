package config

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	key := "a"
	value := "b"
	redis := InitRedis()
	ctx := context.Background()
	result := redis.Set(key, value, ctx, 0)
	if !result {
		t.Fail()
	}

	valueResult := redis.Get(key, ctx)
	if valueResult == nil {
		t.Fail()
	}
	if !strings.EqualFold(fmt.Sprint(valueResult), value) {
		t.Fail()
	}
}
