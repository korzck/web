package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"DB_HOST"`
		Port     int    `yaml:"DB_PORT"`
		Username string `yaml:"DB_USER"`
		Password string `yaml:"DB_PASS"`
		Name     string `yaml:"DB_NAME"`
	} `yaml:"database"`
	Minio struct {
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"minio"`
}

func NewConfig() *Config {
	yamlFile, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		panic(err)
	}

	config := &Config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func DsnFromConf(config *Config) string {
	host := config.Database.Host
	port := config.Database.Port
	user := config.Database.Username
	// pass := config.Database.Password
	dbname := config.Database.Name
	return fmt.Sprintf("host=%v user=%v dbname=%v port=%v sslmode=disable",
		host, user, dbname, port)
}
