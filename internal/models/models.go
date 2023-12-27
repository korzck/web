package models

import "gorm.io/gorm"

type ItemPrototype struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Price    string `json:"price"`
	ImgURL   string `json:"imgurl"`
	URL      string `json:"url"`
	Info     string `json:"info"`
	Type     string `json:"type"`
}
type Item struct {
	gorm.Model
	Id uint `gorm:"primarykey" json:"id"`
	ItemPrototype
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

type UserCreds struct {
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password"`
}

type UserPrototype struct {
	UserCreds
	Name string `json:"name"`
	Tags string `json:"tags"`
}

type User struct {
	gorm.Model
	Id uint `gorm:"primarykey" json:"id"`
	UserPrototype
}

type Order struct {
	gorm.Model
	Id      int    `gorm:"primarykey" json:"id"`
	Status  string `json:"status"`
	UserId  uint   `json:"user_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AdminId uint64 `json:"admin_id"`
}

type OrderItem struct {
	gorm.Model
	OrderId uint
	Order   Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemId  uint
	Item    Item   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment string `json:"comment"`
}
