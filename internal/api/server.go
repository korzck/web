package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Text string
	URL  string
}

func StartServer() {

	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	items := []Item{
		{
			Text: "mashina",
			URL:  "/image/image1.jpg",
		},
		{
			Text: "kvartira",
			URL:  "/image/image2.jpg",
		},
		{
			Text: "zachet",
			URL:  "/image/image3.jpg",
		},
	}

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Items": items,
		})
	})

	r.Static("/image", "./resources")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
