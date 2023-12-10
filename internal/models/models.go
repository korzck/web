package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Id       uint   `gorm:"primarykey" json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Price    string `json:"price"`
	ImgURL   string `json:"imgurl"`
	URL      string `json:"url"`
	Info     string `json:"info"`
	Type     string `json:"type"`
}

type ItemModel struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Price    string `json:"price"`
	ImgURL   string `json:"imgurl"`
	URL      string `json:"url"`
	Info     string `json:"info"`
	Type     string `json:"type"`
}

type User struct {
	gorm.Model
	Id       uint   `gorm:"primarykey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Tags     string `json:"tags"`
}

type Order struct {
	gorm.Model
	Id      int    `gorm:"primarykey" json:"id"`
	Status  string `json:"status"`
	UserId  uint   `json:"user_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment string `json:"comment"`
}

type OrderItem struct {
	gorm.Model
	OrderId uint
	Order   Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemId  uint
	Item    Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
