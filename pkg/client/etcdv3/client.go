package etcdv3

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var Cli *Etcd

// Etcd
type Etcd struct {
	*clientv3.Client
	endpoints   []string
	dialTimeout time.Duration
}

type Option func(*Etcd)

func WithEndpoints(endpoints []string) Option {
	return func(e *Etcd) { e.endpoints = endpoints }
}

func WithDialTimeout(duration time.Duration) Option {
	return func(e *Etcd) { e.dialTimeout = duration * time.Second }
}

func NewClient(opts ...Option) *Etcd {
	var err error
	e := &Etcd{}

	for _, o := range opts {
		o(e)
	}

	e.Client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.endpoints,
		DialTimeout: e.dialTimeout,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return e
}
