package etcdv3

import "time"

type Config struct {
	EndPoints   string
	DialTimeout time.Duration
}
