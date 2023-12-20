package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"web/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Service) GetOrders(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	orders := make([]models.Order, 0)
	ores := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Find(&orders)
	if ores.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": ores.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (s *Service) GetOrder(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	orders := make([]models.Order, 0)
	ores := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("id = ?", orderId).Find(&orders)
	if ores.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": ores.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (s *Service) GetCart(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		// c.JSON(http.StatusOK, gin.H{
		// 	"error": err.Error(),
		// })
		return
	}

	orders := make([]models.Order, 0)
	ores := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("status = 'new'").Find(&orders)
	if ores.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": ores.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (s *Service) Test(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": id,
	})
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
	orderId, ok := c.GetQuery("order_id")

	// jsonData, err := c.GetRawData()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 	return
	// }
	// req := &OrderItemsReq{}
	// err = json.Unmarshal(jsonData, req)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// }
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "no order_id"})
		return
	}
	orderItems := make([]models.OrderItem, 0)
	res := s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", orderId).Find(&orderItems)
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
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *Service) AddItemToOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	jsonData, err := c.GetRawData()
	fmt.Printf("got req %v", string(jsonData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	orderItem := &struct {
		OrderId uint64   `json:"order_id"`
		Items   []uint64 `json:"items"`
	}{}

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
}

func (s *Service) AddItemToCart(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	jsonData, err := c.GetRawData()
	fmt.Printf("got req %v", string(jsonData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	orderItem := &struct {
		OrderId uint64 `json:"order_id"`
		ItemId  uint64 `json:"item_id"`
	}{}

	err = json.Unmarshal(jsonData, orderItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	order := models.Order{}
	s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("status = 'new'").First(&order)
	orderItem.OrderId = uint64(order.Id)
	s.db.DB.Save(&models.OrderItem{
		OrderId: uint(orderItem.OrderId),
		ItemId:  uint(orderItem.ItemId),
	})

	c.JSON(http.StatusOK, orderItem)
}

func (s *Service) DeleteItemFromOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jsonData, err := c.GetRawData()
	fmt.Printf("got data %v", string(jsonData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	orderItem := &struct {
		OrderId uint64 `json:"order_id"`
		Item    uint64 `json:"item_id"`
	}{}

	err = json.Unmarshal(jsonData, orderItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	// for _, v := range orderItem.Items {
	// s.db.DB.Delete()

	// s.db.DB.First(&orderItem, "order_id = ?", )
	foundItem := &models.OrderItem{}
	s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", orderItem.OrderId).Where("item_id = ?", orderItem.Item).First(&foundItem)
	s.db.DB.Delete(foundItem)
	// 		OrderId: uint(orderItem.OrderId),
	// 		ItemId:  uint(v),
	// 	})
	// }

	c.JSON(http.StatusOK, orderItem)
}

func (s *Service) ChangeStatus(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	order := &struct {
		OrderId uint64 `json:"order_id"`
		Status  string `json:"status"`
	}{}

	err = json.Unmarshal(jsonData, order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	foundOrder := &models.Order{}
	s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", order.OrderId).First(&foundOrder)
	foundOrder.Status = order.Status
	tx := s.db.DB.Where("id = ?", order.OrderId).Updates(foundOrder)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *Service) GetOrderInfo(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, _ := c.GetQuery("id")
	foundOrder := &models.Order{}
	s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", id).First(&foundOrder)

	c.JSON(http.StatusOK, foundOrder)
}

func (s *Service) ChangeComment(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	order := &struct {
		Comment string `json:"comment"`
	}{}

	err = json.Unmarshal(jsonData, order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	foundOrder := &models.Order{}
	s.db.DB.Where("deleted_at IS NULL").Where("status = 'new'").Where("user_id = ?", id).First(&foundOrder)
	foundOrder.Comment = order.Comment
	tx := s.db.DB.Where("id = ?", foundOrder.Id).Updates(foundOrder)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
