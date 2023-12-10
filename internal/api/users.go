package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Service) GetOrders(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	orders := make([]models.Order, 0)
	ores := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Find(&orders)
	if ores.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": ores.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

type ItemInOrder struct {
	Id       uint64      `json:"id"`
	Quantity uint64      `json:"quantity"`
	Model    models.Item `json:"model"`
}

type OrderItemsReq struct {
	OrderId uint64 `json:"order_id"`
}

func (s *Service) GetItemsInOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	req := &OrderItemsReq{}
	err = json.Unmarshal(jsonData, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}
	// fmt.Println("got req", req)
	if req.OrderId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "no order_id"})
		return
	}
	orderItems := make([]models.OrderItem, 0)
	res := s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", req.OrderId).Find(&orderItems)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}

	resp := make([]ItemInOrder, 0)
	itemsMap := make(map[uint64]ItemInOrder)
	for _, v := range orderItems {
		q := 0
		if value, ok := itemsMap[uint64(v.ItemId)]; ok {
			q = int(value.Quantity)
		}
		item := &models.Item{}
		s.db.DB.Where("deleted_at IS NULL").Where("id = ?", v.ItemId).Find(&item)
		itemsMap[uint64(v.ItemId)] = ItemInOrder{
			Id:       uint64(item.Id),
			Quantity: uint64(q + 1),
			Model:    *item,
		}
	}
	for _, v := range itemsMap {
		resp = append(resp, v)
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Service) MakeOrder(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	order := &models.Order{}
	err = json.Unmarshal(jsonData, order)
	if err == nil {
		idNum, _ := strconv.ParseInt(id, 10, 64)
		order.UserId = uint(idNum)
		tx := s.db.DB.Save(order)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *Service) AddItemToOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	orderItem := &struct {
		OrderId uint64   `json:"order_id"`
		Items   []uint64 `json:"items"`
	}{}
	// order := &models.OrderItem{}
	err = json.Unmarshal(jsonData, orderItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	for _, v := range orderItem.Items {
		s.db.DB.Save(&models.OrderItem{
			OrderId: uint(orderItem.OrderId),
			ItemId:  uint(v),
		})
	}
	c.JSON(http.StatusOK, orderItem)

	c.JSON(http.StatusOK, gin.H{})
}
