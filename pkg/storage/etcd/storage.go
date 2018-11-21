package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/supergiant/analyze/pkg/storage"
)

type ETCDStorage struct {
	cfg    clientv3.Config
	client *clientv3.Client
}

func (e *ETCDStorage) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	kv := clientv3.NewKV(e.client)

	res, err := kv.Get(ctx, prefix+key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from the etcd")
	}
	if res.Count == 0 {
		return nil, storage.ErrNotFound
	}
	return res.Kvs[0].Value, nil
}

func (e *ETCDStorage) Put(ctx context.Context, prefix string, key string, value []byte) error {
	kv := clientv3.NewKV(e.client)
	_, err := kv.Put(ctx, prefix+key, string(value))
	return errors.Wrap(err, "failed to write to the etcd")
}

func (e *ETCDStorage) Delete(ctx context.Context, prefix string, key string) error {
	_, err := e.client.Delete(ctx, prefix+key, clientv3.WithPrefix())
	return errors.Wrap(err, "failed to read from the etcd")
}

func (e *ETCDStorage) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	result := make([][]byte, 0)
	kv := clientv3.NewKV(e.client)

	r, err := kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return result, errors.Wrap(err, "failed to read from the etcd")
	}
	for _, v := range r.Kvs {
		result = append(result, v.Value)
	}
	return result, nil
}

func NewETCDStorage(cfg clientv3.Config) (storage.Interface, error) {
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the etcd")
	}

	return &ETCDStorage{
		cfg:    cfg,
		client: client,
	}, nil
}

func (e *ETCDStorage) Close() error {
	return e.client.Close()
}
