package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"web/internal"

	"github.com/gin-gonic/gin"
)

func StartServer() error {

	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	items := internal.GetItems()

	for _, v := range items {
		url := fmt.Sprintf("%v", strings.ReplaceAll(v.Title, " ", "-"))
		v.URL = "/" + url
		registerHandler(r, *v)
	}

	r.LoadHTMLGlob("templates/*")
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

		res := make([]*internal.Item, 0)
		// fmt.Println(p1, p2)
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
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		c.HTML(http.StatusOK, "service.html", gin.H{
			"Item": items[id],
		})
	})

	r.Static("/static", "./resources")

	r.Run()

	log.Println("Server down")
	return nil
}

func registerHandler(r *gin.Engine, item internal.Item) {
	url := fmt.Sprintf("%v", strings.ReplaceAll(item.Title, " ", "-"))
	item.URL = "/" + url
	r.GET(item.URL, func(c *gin.Context) {
		c.HTML(http.StatusOK, fmt.Sprintf("%s.html", url), gin.H{
			"Item": item,
		})
	})
}
