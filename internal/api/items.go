package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"web/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	swaggerUtils "github.com/swaggo/swag/example/celler/httputil"
)

// GetItems godoc
// @Summary      Get list of all items
// @Tags         items
// @Param        min       	 query     string  false  "filter by min price"  Format(text)
// @Param        max       	 query     string  false  "filter by max price"  Format(text)
// @Param        title       query     string  false  "filter by title"  Format(text)
// @Param        page    	 query     string  false  "page"  Format(text)
// @Param        material    query     string  false  "filter by material (wood/metal)"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ItemsSwagger
// @Router       /items [get]
func (s *Service) GetItems(c *gin.Context) {
	userId, _, _ := s.getUserRole(c)
	items := make([]models.Item, 0)
	min := c.Request.URL.Query().Get("min")
	max := c.Query("max")
	title := c.Query("title")
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	material := c.Query("material")
	log.Println(min, max, material)
	query := s.db.DB.Where("deleted_at IS NULL")
	if min != "" {
		minNum, _ := strconv.ParseInt(min, 10, 64)
		query.Where(`price >= ?`, minNum)
	}

	if max != "" {
		// log.Println("not empty max")
		maxNum, _ := strconv.ParseInt(max, 10, 64)
		query.Where(`price <= ?`, maxNum)
	}
	if material != "" {
		query.Where(`type = ?`, material)
	}
	if title != "" {
		// q := fmt.Sprintf("title = '%%%v%%'", title)
		query.Where(`title ilike ?`, "%"+title+"%")
	}
	if page == 0 {
		page = 1
	}

	res := query.Order("id ASC").Offset(int(page-1) * 16).Limit(16).Find(&items)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}
	length := struct {
		Count uint64 `json:"count"`
	}{}
	query = s.db.DB.Where("deleted_at IS NULL").Model(&models.Item{})
	if min != "" {
		minNum, _ := strconv.ParseInt(min, 10, 64)
		query.Where(`price >= ?`, minNum)
	}
	if max != "" {
		maxNum, _ := strconv.ParseInt(max, 10, 64)
		query.Where(`price <= ?`, maxNum)
	}
	if material != "" {
		query.Where(`type = ?`, material)
	}
	if title != "" {
		query.Where(`title ilike ?`, "%"+title+"%")
	}
	res = query.Select("count(*)").Find(&length)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}
	// length, _ := res.Rows()
	fmt.Println(length)
	order := models.Order{}
	s.db.DB.Where("user_id = ?", userId).Where("status = ?", "new").First(&order)
	c.JSON(http.StatusOK, &models.ItemsSwagger{
		Items:    items,
		OrderId:  uint64(order.Id),
		Length:   length.Count,
		PageSize: 8,
	})
}

// LoadS3 godoc
// @Summary      Upload s3 file
// @Tags         items
// @Param file formData file true "upload file"
// @Param metadata formData string false "metadata"
// @Accept       mpfd
// @Accept       json
// @Produce      json
// @Success      200  {object} models.ImageSwagger
// @Router       /items/image [post]
func (s *Service) LoadS3(c *gin.Context) {
	// file, err := c.FormFile("file")

	var form models.FormSwagger
	err := c.ShouldBind(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	extension := filepath.Ext(form.File.Filename)
	newFileName := uuid.New().String() + extension
	contentType := form.File.Header["Content-Type"][0]
	buffer, err := form.File.Open()
	if err != nil {
		swaggerUtils.NewError(c, http.StatusBadRequest, err)
		return
	}
	s.minioClient.PutObject(context.Background(), "cnc", newFileName, buffer, form.File.Size, minio.PutObjectOptions{ContentType: contentType})
	reqParams := make(url.Values)
	link, err := s.minioClient.PresignedGetObject(context.Background(), "cnc", newFileName, 7*24*time.Hour, reqParams)
	if link == nil {
		c.JSON(http.StatusInternalServerError, models.ImageSwagger{
			Link:  "",
			Error: err.Error(),
		})
	}
	linkPart := strings.Split(link.String(), "?")
	c.JSON(http.StatusOK, models.ImageSwagger{
		Link:  linkPart[0],
		Error: "",
	})
}

// GetItem godoc
// @Summary      Get item by id
// @Tags         items
// @Param        id    path     string  false  "item id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ItemModel
// @Router       /items/{id} [get]
func (s *Service) GetItem(c *gin.Context) {

	itemId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// res := query.Find(&items)
	item := &models.Item{}
	res := s.db.DB.Where("deleted_at IS NULL").Where("id = ?", itemId).First(&item)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// PostItem godoc
// @Summary      Create item
// @Tags         items
// @Param        itemPrototype body models.ItemPrototype true "Item object"
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ItemPrototype
// @Router       /items/post [post]
func (s *Service) PostItem(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
	}
	item := &models.Item{}
	err = json.Unmarshal(jsonData, item)
	if err == nil {
		tx := s.db.DB.Save(item)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}
}

// Deleteitem godoc
// @Summary      Delete item by id
// @Tags         items
// @Param        id    path     int  true  "item id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /items/{id}/delete [delete]
func (s *Service) DeleteItem(c *gin.Context) {

	itemId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// res := query.Find(&items)
	item := &models.Item{}
	res := s.db.DB.Where("deleted_at IS NULL").Where("id = ?", itemId).Delete(&item)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// PutItem godoc
// @Summary      Change item
// @Tags         items
// @Param        itemPrototype body models.ItemPrototype true "Item object"
// @Param        id    path     int  true  "item id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /items/{id}/put [put]
func (s *Service) PutItem(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
	}
	itemId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item := &models.Item{}
	err = json.Unmarshal(jsonData, item)
	if err == nil {
		tx := s.db.DB.Where("id = ?", itemId).Updates(item)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}
}

// PostItemOrder godoc
// @Summary      Post item to current order
// @Tags         items
// @Param        id    path     int  true  "item id"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200 {object} models.Order
// @Router       /items/{id}/post [post]
func (s *Service) PostItemToOrder(c *gin.Context) {
	userId, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	itemId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	order := &models.Order{}
	tx := s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", userId).Where("status = 'new'").First(&order)
	if tx.Error != nil {
		id, _ := strconv.ParseUint(userId, 10, 64)
		order.UserId = uint(id)
		order.Status = "new"
		txSave := s.db.DB.Save(order)
		if txSave.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
			return
		}
	}

	items := make([]models.OrderItem, 0)
	tx = s.db.DB.Where("deleted_at IS NULL").Where("order_id = ?", order.Id).Where("item_id = ?", itemId).Where("comment != ?", "").First(&items)
	// if tx.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 	return
	// }
	comment := ""
	if len(items) != 0 {
		comment = items[0].Comment
	}
	// orderItem.OrderId = uint64(order.Id)
	s.db.DB.Save(&models.OrderItem{
		OrderId: uint(order.Id),
		ItemId:  uint(itemId),
		Comment: comment,
	})
	tx = s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", userId).Where("status = 'new'").First(&order)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
