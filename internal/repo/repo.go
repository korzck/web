package repo

import (
	"os"

	conf "web/internal/config"
	"web/internal/models"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo() (*Repo, error) {
	yamlFile, err := os.ReadFile("internal/config/config.yaml")
	if err != nil {
		return nil, err
	}

	config := conf.Config{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.Open(conf.DsnFromConf(&config)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repo{
		DB: db,
	}, nil
}

func (r *Repo) GetItems() ([]models.Item, error) {
	items := make([]models.Item, 0)
	res := r.DB.Where("deleted_at IS NULL").Find(&items)
	return items, res.Error
}
