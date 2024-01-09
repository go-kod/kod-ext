package klru

import (
	"context"
	"fmt"

	"dario.cat/mergo"
	"github.com/go-kod/kod-ext/client/kredis"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/samber/lo"
)

type Cache[K comparable, V any] struct {
	cache *lru.Cache[string, V]

	cancel context.CancelFunc
	config *Config
	redis  *kredis.Client
}

type Config struct {
	Size     int
	Channel  string
	UniqueID string
	Redis    kredis.Config
}

func New[K comparable, V any](c *Config) *Cache[K, V] {
	if c.UniqueID == "" {
		panic("UniqueID must be set, which is used to distinguish different cache")
	}

	lo.Must0(mergo.Merge(c, Config{
		Size:    1024,
		Channel: "lru:event:remove",
	}))

	cache := &Cache[K, V]{
		config: c,
		cache:  lo.Must(lru.New[string, V](c.Size)),
	}

	if c.Redis.Addr != "" {
		cache.redis = c.Redis.Build()
		ctx, cancel := context.WithCancel(context.Background())
		cache.cancel = cancel

		go func() {
			pubsub := cache.redis.Subscribe(ctx, c.Channel)
			defer pubsub.Close()

			for {
				select {
				case msg := <-pubsub.Channel():
					cache.cache.Remove(msg.Payload)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return cache
}

func (l *Cache[K, V]) BatchGetOrLoad(ids []K, fn func(missIds []K) (map[K]V, error)) map[K]V {
	res := make(map[K]V)
	missIds := make([]K, 0, len(ids))

	for _, id := range lo.Uniq(ids) {
		if val, ok := l.cache.Get(l.getKey(id)); ok {
			res[id] = val
		} else {
			missIds = append(missIds, id)
		}
	}

	if len(missIds) == 0 || fn == nil {
		return res
	}

	dataMap, err := fn(missIds)
	if err != nil || dataMap == nil {
		return res
	}

	for id, val := range dataMap {
		l.cache.Add(l.getKey(id), val)
		res[id] = val
	}

	return res
}

func (l *Cache[K, V]) RemoveAll(ctx context.Context, id K) error {
	if l.redis != nil {
		return l.redis.Publish(ctx, l.config.Channel, l.getKey(id)).Err()
	}

	return nil
}

func (l *Cache[K, V]) getKey(id K) string {
	return fmt.Sprintf("%s:%v", l.config.UniqueID, id)
}
