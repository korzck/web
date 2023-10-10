package main

import (
	"os"

	conf "web/internal/config"
	"web/internal/models"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	yamlFile, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		panic(err)
	}

	config := conf.Config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(conf.DsnFromConf(&config)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Order{}, &models.Item{}, &models.OrderItem{})
	if err != nil {
		panic("cant migrate db")
	}
}
