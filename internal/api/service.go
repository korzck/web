package api

import (
	"log"
	"net/http"

	mClient "web/internal/minio"
	"web/internal/redis"
	"web/internal/repo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	db          *repo.Repo
	minioClient *mClient.MinioClient
	redisRepo   *redis.RedisRepo
}

func NewService(
	db *repo.Repo,
	minioClient *mClient.MinioClient,
	redisRepo *redis.RedisRepo,
) *Service {
	return &Service{
		db:          db,
		minioClient: minioClient,
		redisRepo:   redisRepo,
	}
}

type Empty struct{}

func (s *Service) Run() error {

	log.Println("Server start up")

	r := gin.Default()
	r.Use(CORSMiddleware())

	// db, _ := repo.NewRepo()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/s3/upload", s.LoadS3)

	r.GET("/items", s.GetItems)

	r.POST("/signup", s.SignUp)
	r.POST("/login", s.Login)
	r.POST("/logout", s.Logout)
	r.Use(s.UserAuth).POST("/validate", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	r.Use(s.AdminAuth).POST("/validate_admin", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/orders", s.GetOrders)
	r.POST("/orders", s.MakeOrder)

	r.GET("/orderitems", s.GetItemsInOrder)
	r.POST("/orderitems", s.AddItemToOrder)

	// admin := r.Group("/admin")
	// {

	// 	admin.POST("/item", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
	// 		}
	// 		item := &models.Item{}
	// 		err = json.Unmarshal(jsonData, item)
	// 		if err == nil {
	// 			tx := db.DB.Save(item)
	// 			if tx.Error != nil {
	// 				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 				return
	// 			}
	// 			c.JSON(http.StatusOK, item)
	// 		} else {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 		}
	// 	})

	// 	admin.PUT("/item", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		item := &models.Item{}
	// 		err = json.Unmarshal(jsonData, item)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		if item.Id == 0 {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
	// 			return
	// 		}
	// 		tx := db.DB.Where("id = ?", item.Id).Updates(item)
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, item)

	// 	})

	// 	admin.DELETE("/item", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		id := struct {
	// 			Id uint64 `json:"id"`
	// 		}{}
	// 		err = json.Unmarshal(jsonData, &id)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		tx := db.DB.Delete(&models.Item{
	// 			Id: uint(id.Id),
	// 		})
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
	// 	})

	// 	// =================================================================

	// 	admin.POST("/user", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		user := &models.User{}
	// 		err = json.Unmarshal(jsonData, user)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 		}
	// 		tx := db.DB.Save(user)
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, user)
	// 	})

	// 	admin.PUT("/user", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		user := &models.User{}
	// 		err = json.Unmarshal(jsonData, user)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		if user.Id == 0 {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
	// 			return
	// 		}
	// 		tx := db.DB.Where("id = ?", user.Id).Updates(user)
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, user)

	// 	})

	// 	admin.DELETE("/user", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		id := struct {
	// 			Id uint64 `json:"id"`
	// 		}{}
	// 		err = json.Unmarshal(jsonData, &id)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		tx := db.DB.Delete(&models.User{
	// 			Id: uint(id.Id),
	// 		})
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
	// 	})

	// 	// =================================================================

	// 	admin.POST("/order", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		order := &models.Order{}
	// 		err = json.Unmarshal(jsonData, order)
	// 		if err == nil {
	// 			tx := db.DB.Save(order)
	// 			if tx.Error != nil {
	// 				c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 				return
	// 			}
	// 			c.JSON(http.StatusOK, order)
	// 			return
	// 		} else {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 		}
	// 	})

	// 	admin.PUT("/order", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		order := &models.Order{}
	// 		err = json.Unmarshal(jsonData, order)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		if order.Id == 0 {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
	// 			return
	// 		}
	// 		tx := db.DB.Where("id = ?", order.Id).Updates(order)
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, order)

	// 	})

	// 	admin.DELETE("/order", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		id := struct {
	// 			Id uint64 `json:"id"`
	// 		}{}
	// 		err = json.Unmarshal(jsonData, &id)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		tx := db.DB.Delete(&models.Order{
	// 			Id: int(id.Id),
	// 		})
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
	// 	})

	// 	// =================================================================

	// 	admin.POST("/orderitem", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		orderItem := &struct {
	// 			OrderId uint64   `json:"order_id"`
	// 			Items   []uint64 `json:"items"`
	// 		}{}
	// 		// order := &models.OrderItem{}
	// 		err = json.Unmarshal(jsonData, orderItem)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		for _, v := range orderItem.Items {
	// 			db.DB.Save(&models.OrderItem{
	// 				OrderId: uint(orderItem.OrderId),
	// 				ItemId:  uint(v),
	// 			})
	// 		}
	// 		c.JSON(http.StatusOK, orderItem)
	// 	})

	// 	admin.PUT("/orderitem", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		order := &models.Order{}
	// 		err = json.Unmarshal(jsonData, order)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		if order.Id == 0 {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": "no id provided"})
	// 			return
	// 		}
	// 		tx := db.DB.Where("id = ?", order.Id).Updates(order)
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, order)

	// 	})

	// 	admin.DELETE("/orderitem", func(c *gin.Context) {
	// 		jsonData, err := c.GetRawData()
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		id := struct {
	// 			Order uint64 `json:"order_id"`
	// 			Item  uint64 `json:"item_id"`
	// 		}{}
	// 		err = json.Unmarshal(jsonData, &id)
	// 		if err != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	// 			return
	// 		}
	// 		tx := db.DB.Where("order_id = ?", id.Order).Where("item_id = ?", id.Item).Delete(&models.OrderItem{})
	// 		if tx.Error != nil {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 			return
	// 		}
	// 		c.JSON(http.StatusOK, gin.H{"error: ": tx.Error})
	// 	})

	// }

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
