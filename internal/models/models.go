package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Id       uint `gorm:"primarykey"`
	Title    string
	Subtitle string
	Price    string
	ImgURL   string
	URL      string
	Info     string
	Type     string
}

type User struct {
	gorm.Model
	Id    uint   `gorm:"primarykey"`
	Email string `gorm:"unique;not null"`
	Name  string
}

type Order struct {
	gorm.Model
	Id     int `gorm:"primarykey"`
	UserId uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderItem struct {
	gorm.Model
	OrderId uint
	Order   Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemId  uint
	Item    Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
