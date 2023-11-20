package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	// "web/internal"

	mClient "web/internal/minio"
	"web/internal/models"
	"web/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func StartServer() error {

	log.Println("Server start up")

	r := gin.Default()

	db, _ := repo.NewRepo()

	items, _ := db.GetItems()
	// r.LoadHTMLGlob("templates/*")
	var err error
	r.GET("/services", func(c *gin.Context) {
		filter := c.Query("filter")
		priceFrom := c.Query("pricefrom")
		priceUpTo := c.Query("priceupto")
		if filter == "" {
			filter = "all"
		}
		var p1, p2 int64
		if priceFrom == "" {
			p1 = 0
		} else {
			p1, err = strconv.ParseInt(priceFrom, 10, 64)
			if err != nil {
				return
			}
		}
		if priceUpTo == "" {
			p2 = int64(^uint32(0))
		} else {
			p2, err = strconv.ParseInt(priceUpTo, 10, 64)
			if err != nil {
				return
			}
		}

		res := make([]models.Item, 0)
		for _, v := range items {
			p, _ := strconv.ParseInt(v.Price, 10, 64)
			if p < p1 || p >= p2 {
				continue
			}
			if filter == "all" {
				res = append(res, v)
			} else if v.Type == filter {
				res = append(res, v)
			}
		}

		c.HTML(http.StatusOK, "services.html", gin.H{
			"Items":     res,
			"Filter":    filter,
			"PriceFrom": priceFrom,
			"PriceUpTo": priceUpTo,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/contacts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contacts.html", gin.H{})
	})

	r.GET("/services/:id", func(c *gin.Context) {
		items, _ = db.GetItems()
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		c.HTML(http.StatusOK, "service.html", gin.H{
			"Item": items[id-1],
		})
	})

	admin := r.Group("/admin")
	{
		admin.GET("/additem", func(c *gin.Context) {
			items, _ = db.GetItems()
			c.HTML(http.StatusOK, "additem.html", gin.H{
				"Items": items,
			})
		})
		admin.POST("/additem", func(c *gin.Context) {
			items, _ = db.GetItems()
			c.HTML(http.StatusOK, "additem.html", gin.H{
				"Items": items,
			})
			title := c.PostForm("title")
			subtitle := c.PostForm("subtitle")
			price := c.PostForm("price")
			url := c.PostForm("url")
			info := c.PostForm("info")
			typeText := c.PostForm("type")
			if title == "" {
				return
			}
			db.DB.Save(&models.Item{
				Title:    title,
				Subtitle: subtitle,
				Price:    price,
				ImgURL:   url,
				Info:     info,
				Type:     typeText,
			})

		})
		admin.POST("/deleteitem/:id", func(c *gin.Context) {
			// c.HTML(http.StatusOK, "additem.html", gin.H{
			// 	"Items": items,
			// })
			item := models.Item{}
			// id := c.PostForm("title")
			// id := c.PostForm("id")
			id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
			db.DB.Delete(&item, id)
			c.Redirect(304, "/admin/additem")
		})

	}

	r.Static("/static", "./resources")

	r.Run()

	log.Println("Server down")
	return nil
}

func Run() error {

	log.Println("Server start up")

	r := gin.Default()
	r.Use(CORSMiddleware())

	db, _ := repo.NewRepo()

	minioClient := mClient.NewMinioClient()

	// items, _ := db.GetItems()
	// var err error

	r.POST("/s3/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// The file cannot be received.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
			return
		}

		// Retrieve file information
		extension := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + extension
		// filePath := "/files/" + newFileName
		contentType := file.Header["Content-Type"][0]
		buffer, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		minioClient.PutObject(context.Background(), "cnc", newFileName, buffer, file.Size, minio.PutObjectOptions{ContentType: contentType})
		reqParams := make(url.Values)
		link, err := minioClient.PresignedGetObject(context.Background(), "cnc", newFileName, 7*24*time.Hour, reqParams)
		if link == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"link": "",
				"err":  err,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"link": link.String(),
			"err":  err,
		})
	})

	r.GET("/items", func(c *gin.Context) {
		items := make([]models.Item, 0)
		res := db.DB.Where("deleted_at IS NULL").Find(&items)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	r.GET("/orders", func(c *gin.Context) {
		jsonData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		id := struct {
			UserId uint64 `json:"user_id"`
		}{}
		err = json.Unmarshal(jsonData, &id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		orders := make([]models.Order, 0)
		res := db.DB.Where("deleted_at IS NULL").Where("user_id = ?", id.UserId).Find(&orders)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	r.GET("/orderitems", func(c *gin.Context) {
		jsonData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		id := struct {
			OrderId uint64 `json:"order_id"`
		}{}
		err = json.Unmarshal(jsonData, &id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}
		orderItems := make([]models.OrderItem, 0)
		res := db.DB.Where("deleted_at IS NULL").Where("order_id = ?", id.OrderId).Find(&orderItems)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, orderItems)
	})

	admin := r.Group("/admin")
	{

		admin.POST("/item", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
			}
			item := &models.Item{}
			err = json.Unmarshal(jsonData, item)
			if err == nil {
				tx := db.DB.Save(item)
				if tx.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
					return
				}
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			}
		})

		admin.PUT("/item", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			item := &models.Item{}
			err = json.Unmarshal(jsonData, item)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			if item.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
				return
			}
			tx := db.DB.Where("id = ?", item.Id).Updates(item)
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, item)

		})

		admin.DELETE("/item", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			id := struct {
				Id uint64 `json:"id"`
			}{}
			err = json.Unmarshal(jsonData, &id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			tx := db.DB.Delete(&models.Item{
				Id: uint(id.Id),
			})
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
		})

		// =================================================================

		admin.POST("/user", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			user := &models.User{}
			err = json.Unmarshal(jsonData, user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			}
			tx := db.DB.Save(user)
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, user)
		})

		admin.PUT("/user", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			user := &models.User{}
			err = json.Unmarshal(jsonData, user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			if user.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
				return
			}
			tx := db.DB.Where("id = ?", user.Id).Updates(user)
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, user)

		})

		admin.DELETE("/user", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			id := struct {
				Id uint64 `json:"id"`
			}{}
			err = json.Unmarshal(jsonData, &id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			tx := db.DB.Delete(&models.User{
				Id: uint(id.Id),
			})
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
		})

		// =================================================================

		admin.POST("/order", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			order := &models.Order{}
			err = json.Unmarshal(jsonData, order)
			if err == nil {
				tx := db.DB.Save(order)
				if tx.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
					return
				}
				c.JSON(http.StatusOK, order)
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			}
		})

		admin.PUT("/order", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			order := &models.Order{}
			err = json.Unmarshal(jsonData, order)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			if order.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
				return
			}
			tx := db.DB.Where("id = ?", order.Id).Updates(order)
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, order)

		})

		admin.DELETE("/order", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			id := struct {
				Id uint64 `json:"id"`
			}{}
			err = json.Unmarshal(jsonData, &id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			tx := db.DB.Delete(&models.Order{
				Id: int(id.Id),
			})
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
		})

		// =================================================================

		admin.POST("/orderitem", func(c *gin.Context) {
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
				db.DB.Save(&models.OrderItem{
					OrderId: uint(orderItem.OrderId),
					ItemId:  uint(v),
				})
			}
			c.JSON(http.StatusOK, orderItem)
		})

		admin.PUT("/orderitem", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			order := &models.Order{}
			err = json.Unmarshal(jsonData, order)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			if order.Id == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
				return
			}
			tx := db.DB.Where("id = ?", order.Id).Updates(order)
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, order)

		})

		admin.DELETE("/orderitem", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			id := struct {
				Order uint64 `json:"order_id"`
				Item  uint64 `json:"item_id"`
			}{}
			err = json.Unmarshal(jsonData, &id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
				return
			}
			tx := db.DB.Where("order_id = ?", id.Order).Where("item_id = ?", id.Item).Delete(&models.OrderItem{})
			if tx.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
		})

	}

	r.Run()

	log.Println("Server down")
	return nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
