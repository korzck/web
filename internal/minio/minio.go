package minio

import (
	"log"
	"os"

	conf "web/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"gopkg.in/yaml.v2"
)

type MinioClient struct {
	*minio.Client
}

func NewMinioClient() *MinioClient {

	yamlFile, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	config := conf.Config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalln(err)
	}
	useSSL := false

	minioClient, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.User, config.Minio.Pass, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &MinioClient{
		minioClient,
	}
}
