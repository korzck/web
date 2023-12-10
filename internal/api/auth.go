package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"web/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SignUp(c *gin.Context) {
	user := &models.User{}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user.Password = string(hash)
	user.Tags = ""
	tx := s.db.DB.Save(user)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error saving user: " + tx.Error.Error(),
		})
		return
	}
}

func (s *Service) Login(c *gin.Context) {
	user := &models.User{}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	foundUser := &models.User{}
	s.db.DB.First(&foundUser, "email = ?", user.Email)

	if foundUser.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing user",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password or email",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "admin",
		Subject:   fmt.Sprintf("%v", foundUser.Id),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token generate",
		})
		return
	}

	strToken, err := token.SignedString([]byte("test123"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant create jwt token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth", strToken, 3600*24*30, "", "", false, true)
	s.redisRepo.RedisClient.Del(context.Background(), fmt.Sprintf("%v", foundUser.Id))
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Service) getUserRole(c *gin.Context) (string, string, error) {
	cookie, err := c.Cookie("auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte("test123"), nil
	})
	id := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Printf("user id is %v", claims["sub"])
		id = claims["sub"].(string)
	} else {
		return "", "", err
	}
	user := &models.User{}
	s.db.DB.First(user, "id = ?", id)

	return id, user.Tags, nil
}

func (s *Service) AdminAuth(c *gin.Context) {
	id, role, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	if s.isLogout(id) {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
	c.Next()
}

func (s *Service) UserAuth(c *gin.Context) {
	id, role, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	if s.isLogout(id) {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if role != "admin" && role != "user" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

func (s *Service) Logout(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	s.redisRepo.RedisClient.Set(context.Background(), id, "logout", 0)
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Service) isLogout(id string) bool {
	_, err := s.redisRepo.RedisClient.Get(context.Background(), id).Result()
	return err == nil
}
