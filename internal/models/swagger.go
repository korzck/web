package models

import "mime/multipart"

type ItemsSwagger struct {
	Items    []Item `json:"items"`
	OrderId  uint64 `json:"order_id"`
	Length   uint64 `json:"length"`
	PageSize uint64 `json:"page_size"`
}

type UserSwagger struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Id    uint   `json:"id"`
	Tags  string `json:"tags"`
	Order uint64 `json:"order"`
}

type FormSwagger struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type OrderSwagger struct {
	Status string               `json:"status"`
	Id     uint                 `json:"id"`
	UserId uint                 `json:"user_id"`
	Email  string               `json:"email"`
	Items  []ItemInOrderSwagger `json:"items"`
}

type ItemInOrderSwagger struct {
	Id       uint64 `json:"id"`
	Quantity uint64 `json:"quantity"`
	Comment  string `json:"comment"`
	Item     Item   `json:"item"`
}

type OrderStatusSwagger struct {
	Status string `json:"status"`
}

type ItemCommentSwagger struct {
	ItemId  uint64 `json:"item_id"`
	Comment string `json:"comment"`
}

type ImageSwagger struct {
	Link  string `json:"link"`
	Error string `json:"error"`
}
