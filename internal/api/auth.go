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

// Sign godoc
// @Summary      Sign up
// @Tags         auth
// @Param        userPrototype body models.UserPrototype true "User object"
// @Accept       json
// @Produce      json
// @Success      200  {int}  id
// @Router       /signup [post]
func (s *Service) SignUp(c *gin.Context) {
	user := &models.User{}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	foundUser := &models.User{}
	s.db.DB.First(&foundUser, "email = ?", user.Email)

	if foundUser.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email already exists",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user.Password = string(hash)
	// user.Tags = ""
	tx := s.db.DB.Save(user)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error saving user: " + tx.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": user.Id,
	})
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Param        userCreds body models.UserCreds true "User object"
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.UserSwagger
// @Router       /login [post]
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
	res := s.redisRepo.RedisClient.Del(context.Background(), fmt.Sprintf("%v", foundUser.Id))
	if res.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant delete from logout",
		})
		return
	}

	// order := &models.Order{}
	// s.db.DB.Where("deleted_at IS NULL").Where("user_id = ?", foundUser.Id).Where("status = 'new'").First(&order)
	// if order.Status != "new" {
	// 	order.UserId = foundUser.Id
	// 	order.Status = "new"
	// 	tx := s.db.DB.Save(order)
	// 	if tx.Error != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error: ": tx.Error.Error()})
	// 		return
	// 	}
	// }

	c.JSON(http.StatusOK, &models.UserSwagger{
		Id:    foundUser.Id,
		Name:  foundUser.Name,
		Email: foundUser.Email,
		Tags:  foundUser.Tags,
	})
}

func (s *Service) getUserRole(c *gin.Context) (string, string, error) {
	cookie, err := c.Cookie("auth")
	if err != nil {
		// fmt.Printf("got err %v", err)
		return "", "", err
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte("test123"), nil
	})

	id := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
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
		c.AbortWithStatus(http.StatusForbidden)
	}
	if role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)

	}
	c.Next()
}

func (s *Service) UserAuth(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if s.isLogout(id) {
		c.AbortWithStatus(http.StatusForbidden)
	}
	c.Next()
}

// Logout godoc
// @Summary      Logout
// @Tags         auth
// @Param 		 Cookie header string  false "token"     default(token=xxx)
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /logout [post]
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

// Validate godoc
// @Summary      validate auth
// @Tags         auth
// @Param 		 Cookie header string  false "token"     default(token=xxx)
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.UserSwagger
// @Router       /validate [post]
func (s *Service) Validate(c *gin.Context) {
	id, _, err := s.getUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	foundUser := &models.User{}
	foundOrder := &models.Order{}
	s.db.DB.First(&foundUser, "id = ?", id)
	s.db.DB.Where("user_id = ?", foundUser.Id).Where("status = ?", "new").First(&foundOrder)

	c.JSON(http.StatusOK, &models.UserSwagger{
		Email: foundUser.Email,
		Name:  foundUser.Name,
		Id:    foundUser.Id,
		Tags:  foundUser.Tags,
		Order: uint64(foundOrder.Id),
	})

}
