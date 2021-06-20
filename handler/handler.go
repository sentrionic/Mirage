package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/mirage/handler/middleware"
	"github.com/sentrionic/mirage/model"
	"github.com/sentrionic/mirage/model/apperrors"
	"time"
)

type Handler struct {
	UserService  model.UserService
	MaxBodyBytes int64
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	}

	// Create an account group
	ag := c.R.Group("v1/accounts")

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)

	ag.Use(middleware.AuthUser())

	ag.GET("", h.Current)
	ag.PUT("", h.EditAccount)
	ag.POST("/logout", h.Logout)
}

// setUserSession saves the users ID in the session
func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
}
