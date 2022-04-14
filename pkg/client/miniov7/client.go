package miniov7

import (
	"gin-base/pkg/logging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioCli *Client

type Config struct {
	Endpoint string
	User     string
	Password string
	UseSSl   bool
}

var logger = logging.GetLogger("miniov7")

type Client struct {
	*minio.Client
}

func NewClient(config *Config) *Client {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.User, config.Password, ""),
		Secure: config.UseSSl,
	})
	if err != nil {
		panic(err)
	}

	return &Client{minioClient}
}
