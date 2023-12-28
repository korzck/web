package redis

import (
	"context"
	"fmt"
	"os"

	conf "web/internal/config"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
)

type RedisRepo struct {
	RedisClient *redis.Client
}

func New() (*RedisRepo, error) {
	yamlFile, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		return nil, err
	}

	config := conf.Config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: "",
		DB:       0,
	})

	// Ping the Redis server to check the connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}
	fmt.Println("Connected to Redis")
	return &RedisRepo{
		RedisClient: client,
	}, nil
}
