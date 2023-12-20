package api

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
	"web/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/swaggo/swag/example/celler/httputil"
)

// LoadS3 godoc
// @Summary      Upload s3 file
// @Tags         s3
// @Param file formData file true "upload file"
// @Param metadata formData string false "metadata"
// @Accept      mpfd
// @Accept       json
// @Produce      json
// @Success      200  {object}  Empty
// @Router       /s3/upload [post]
func (s *Service) LoadS3(c *gin.Context) {
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	contentType := file.Header["Content-Type"][0]
	buffer, err := file.Open()
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	s.minioClient.PutObject(context.Background(), "cnc", newFileName, buffer, file.Size, minio.PutObjectOptions{ContentType: contentType})
	reqParams := make(url.Values)
	link, err := s.minioClient.PresignedGetObject(context.Background(), "cnc", newFileName, 7*24*time.Hour, reqParams)
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
}

// GetItems godoc
// @Summary      Get list of all items
// @Tags         items
// @Param        min    query     string  false  "filter by min price"  Format(text)
// @Param        max    query     string  false  "filter by max price"  Format(text)
// @Param        material    query     string  false  "filter by material (wood/metal)"  Format(text)
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.ItemModel
// @Router       /items [get]
func (s *Service) GetItems(c *gin.Context) {
	items := make([]models.Item, 0)
	min := c.Request.URL.Query().Get("min")
	max := c.Query("max")
	material := c.Query("material")
	log.Println(min, max, material)
	query := s.db.DB.Where("deleted_at IS NULL")
	if min != "" {
		minNum, _ := strconv.ParseInt(min, 10, 64)
		query.Where(`NULLIF(price, '')::int >= ?`, minNum)
	}

	if max != "" {
		// log.Println("not empty max")
		maxNum, _ := strconv.ParseInt(max, 10, 64)
		query.Where(`NULLIF(price, '')::int <= ?`, maxNum)
	}
	if material != "" {
		query.Where(`type = ?`, material)
	}
	res := query.Find(&items)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": res.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
