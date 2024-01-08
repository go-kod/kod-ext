package klru

import (
	"context"
	"testing"
	"time"

	"github.com/go-kod/kod-ext/client/kredis"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	l1 := New[string, string](&Config{
		UniqueID: "TestCache",
		Channel:  "TestCache",
		Redis: kredis.Config{
			Addr: "localhost:6379",
		},
	})
	res := l1.BatchGetOrLoad([]string{"1"}, func(missKeys []string) (map[string]string, error) {
		return map[string]string{"1": "1"}, nil
	})

	assert.Equal(t, res["1"], "1")

	res = l1.BatchGetOrLoad([]string{"1"}, nil)

	assert.Equal(t, "1", res["1"])

	l2 := New[string, string](&Config{
		UniqueID: "TestCache",
		Channel:  "TestCache",
		Redis: kredis.Config{
			Addr: "localhost:6379",
		},
	})

	res = l1.BatchGetOrLoad([]string{"1"}, nil)

	assert.Equal(t, res["1"], "1")

	l2.RemoveAll(context.Background(), "1")
	time.Sleep(2 * time.Second)

	res = l1.BatchGetOrLoad([]string{"1"}, nil)

	assert.NotEqual(t, res["1"], "1")
}
