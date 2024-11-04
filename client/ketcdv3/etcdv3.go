package ketcdv3

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Config struct {
	Endpoints []string      `default:"-"`
	Timeout   time.Duration `default:"3s"`
}

func (r Config) Build(_ context.Context) (*clientv3.Client, error) {
	if len(r.Endpoints) == 0 {
		return nil, fmt.Errorf("no etcd endpoints provided")
	}

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   r.Endpoints,
		DialTimeout: r.Timeout,
		DialOptions: []grpc.DialOption{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return etcd, nil
}
