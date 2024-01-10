package api

import (
	"encoding/json"
	"log"
	"net/http"

	mClient "web/internal/minio"
	"web/internal/models"
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

	r.POST("/items/image", s.LoadS3)

	r.GET("/items", s.GetItems)
	r.GET("/items/:id", s.GetItem)
	r.POST("/items/post", s.PostItem)
	r.DELETE("/items/:id/delete", s.DeleteItem)
	r.PUT("/items/:id/put", s.PutItem)
	r.POST("/items/:id/post", s.PostItemToOrder)

	r.POST("/signup", s.SignUp)
	r.POST("/login", s.Login)
	r.POST("/logout", s.Logout)

	r.GET("/orders", s.GetOrders)
	r.GET("/orders/:id", s.GetItemsInOrder)
	r.PUT("/orders/:id/approve", s.PutOrderStatus)
	r.PUT("/orders/make", s.ConfirmOrder)
	r.DELETE("/orders/delete", s.DeleteOrder)
	r.DELETE("orders/items/:id", s.DeleteItemFromOrder)
	r.PUT("orders/:id/comment", s.AddItemComment)

	// r.GET("/orders/current", s.GetCart)

	// r.GET("/test", s.Test)
	// r.PUT("/orders", s.ChangeStatus)
	// r.POST("/orders", s.MakeOrder)
	// r.POST("/cartcomment", s.ChangeComment)

	// r.GET("/orderitems", s.GetItemsInOrder)
	// r.POST("/orderitems", s.AddItemToOrder)
	// r.POST("/cartadditem", s.AddItemToCart)
	// r.DELETE("/orderitems", s.DeleteItemFromOrder)

	r.Use(s.UserAuth).POST("/validate", s.Validate)
	r.Use(s.AdminAuth).POST("/validate_admin", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	admin := r.Group("/admin")
	{

		admin.Use(s.AdminAuth).POST("/item", func(c *gin.Context) {
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
		})

		admin.Use(s.AdminAuth).PUT("/item", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
			}
			item := &models.Item{}
			err = json.Unmarshal(jsonData, item)
			if err == nil {
				tx := s.db.DB.Where("id = ?", item.Id).Updates(item)
				if tx.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
					return
				}
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			}
		})
		admin.Use(s.AdminAuth).DELETE("/item", func(c *gin.Context) {
			jsonData, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error: ": err.Error()})
			}
			item := &models.Item{}
			err = json.Unmarshal(jsonData, item)
			if err == nil {
				tx := s.db.DB.Where("id = ?", item.Id).First(&item)
				if tx.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
					return
				}
				s.db.DB.Delete(item)
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			}
		})
	}

	r.Run("10.0.0.21:8080")

	log.Println("Server down")
	return nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
