package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"web/internal/models"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary      Get list of all orders
// @Tags         orders
// @Param        min_date    query     string  false  "min date"  Format(text)
// @Param        max_date    query     string  false  "max date"  Format(text)
// @Param        status      query     string  false  "order status"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Order
// @Router       /orders [get]
func (s *Service) GetOrders(c *gin.Context) {
	id, tag, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	if tag == "admin" {
		adminOrders := make([]models.Order, 0)
		max := c.Query("max_date")
		min := c.Query("min_date")
		status := c.Query("status")
		tx := s.db.DB.Where("deleted_at IS NULL")
		if max != "" {
			date, _ := time.Parse("2006-01-02", max)
			tx = tx.Where("created_at::date <?", date)
		}
		if min != "" {
			date, _ := time.Parse("2006-01-02", min)
			tx = tx.Where("created_at::date  >=?", date)
		}
		if status != "all" {
			tx = tx.Where("status = ?", status)
		} else {
			tx = tx.Where("status != ?", "new")
		}

		tx = tx.Find(&adminOrders)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, adminOrders)
		return
	}

	orders := make([]models.Order, 0)
	ores := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("status != ?", "new").Find(&orders)
	if ores.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": ores.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// PutOrderStatus godoc
// @Summary      Approve or decline order
// @Tags         orders
// @Param        status body models.OrderStatusSwagger true "Order status"
// @Param        id    path     string  true  "order id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/{id}/approve [put]
func (s *Service) PutOrderStatus(c *gin.Context) {
	id, tag, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	if tag != "admin" {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	status := &struct {
		Status string `json:"status"`
	}{}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	err = json.Unmarshal(jsonData, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("id = ?", orderId).First(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}
	order.Status = status.Status
	order.AdminId, _ = strconv.ParseUint(id, 10, 64)
	tx = s.db.DB.Where("id = ?", orderId).Updates(order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// ConfirmOrder godoc
// @Summary      Confirm current order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/make [put]
func (s *Service) ConfirmOrder(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("status = ?", "new").First(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}
	order.Status = "pending"
	order.CreatedAt = time.Now()
	tx = s.db.DB.Where("id = ?", order.Id).Updates(order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteOrder godoc
// @Summary      Delete current order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/delete [delete]
func (s *Service) DeleteOrder(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id).Where("status = ?", "new").Delete(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteItemFromOrder godoc
// @Summary      Delete item from current order
// @Tags         orders
// @Param        id    path     string  true  "item id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200 {object} models.OrderSwagger
// @Router       /orders/items/{id} [delete]
func (s *Service) DeleteItemFromOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	itemId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("status = ?", "new").First(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	foundItem := &models.OrderItem{}
	tx = s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", order.Id).Where("item_id = ?", itemId).First(&foundItem)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}
	tx = s.db.DB.Delete(foundItem)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	orderItems := make([]models.OrderItem, 0)
	tx = s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", order.Id).Find(&orderItems)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	resp := make([]models.ItemInOrderSwagger, 0)
	itemsMap := make(map[uint64]models.ItemInOrderSwagger)
	for _, v := range orderItems {
		q := 0
		if value, ok := itemsMap[uint64(v.ItemId)]; ok {
			q = int(value.Quantity)
		}
		item := &models.Item{}
		s.db.DB.Where("deleted_at IS NULL").Where("id = ?", v.ItemId).Find(&item)
		itemsMap[uint64(v.ItemId)] = models.ItemInOrderSwagger{
			Id:       uint64(item.Id),
			Quantity: uint64(q + 1),
			Item:     *item,
		}
	}
	for _, v := range itemsMap {
		resp = append(resp, v)
	}

	if len(resp) == 0 {
		tx = s.db.DB.Where("deleted_at IS NULL").Where("id = ?", order.Id).Delete(&order)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, &models.OrderSwagger{
			Id: 0,
		})
	}

	c.JSON(http.StatusOK, &models.OrderSwagger{
		Id:     uint(order.Id),
		Status: order.Status,
		UserId: order.UserId,
		Items:  resp,
	})
}

// DeleteItemFromOrder godoc
// @Summary      Delete item from current order
// @Tags         orders
// @Param        id    path     string  true  "item id"  Format(text)
// @Param        comment body models.ItemCommentSwagger true "Item comment"
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/{id}/comment [put]
func (s *Service) AddItemComment(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	orderItems := make([]models.OrderItem, 0)
	res := s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", orderId).Find(&orderItems)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}

	item := &models.ItemCommentSwagger{}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error: ": err.Error(),
		})
		return
	}
	err = json.Unmarshal(jsonData, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}

	for i := range orderItems {
		if orderItems[i].ItemId == uint(item.ItemId) {
			orderItems[i].Comment = item.Comment
		}
	}
	tx := s.db.DB.Save(orderItems)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetOrder godoc
// @Summary      Get order by id
// @Tags         orders
// @Param        id    path     string  true  "order id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.OrderSwagger
// @Router       /orders/{id} [get]
func (s *Service) GetItemsInOrder(c *gin.Context) {
	_, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error: ": "no order_id"})
	// 	return
	// }
	orderItems := make([]models.OrderItem, 0)
	res := s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", orderId).Find(&orderItems)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("id = ?", orderId).First(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	resp := make([]models.ItemInOrderSwagger, 0)
	itemsMap := make(map[uint64]models.ItemInOrderSwagger)
	for _, v := range orderItems {
		q := 0
		if value, ok := itemsMap[uint64(v.ItemId)]; ok {
			q = int(value.Quantity)
		}
		item := &models.Item{}
		s.db.DB.Where("deleted_at IS NULL").Where("id = ?", v.ItemId).Find(&item)
		itemsMap[uint64(v.ItemId)] = models.ItemInOrderSwagger{
			Id:       uint64(item.Id),
			Quantity: uint64(q + 1),
			Comment:  v.Comment,
			Item:     *item,
		}
	}
	for _, v := range itemsMap {
		resp = append(resp, v)
	}

	c.JSON(http.StatusOK, &models.OrderSwagger{
		Id:     uint(orderId),
		Status: order.Status,
		UserId: order.UserId,
		Items:  resp,
	})
}
