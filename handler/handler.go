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
	PostService  model.PostService
	MaxBodyBytes int64
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	PostService     model.PostService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		PostService:  c.PostService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	}

	// Account group
	ag := c.R.Group("v1/accounts")

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)

	ag.Use(middleware.AuthUser())

	ag.GET("", h.Current)
	ag.PUT("", h.EditAccount)
	ag.POST("/logout", h.Logout)

	// User group
	ug := c.R.Group("v1/profiles")
	ug.GET("/:username", h.GetProfile)

	ug.Use(middleware.AuthUser())
	ug.POST("/:username/follow", h.ToggleFollow)

	// Post group
	pg := c.R.Group("v1/posts")

	pg.Use(middleware.AuthUser())

	pg.POST("", h.CreatePost)
}

// setUserSession saves the users ID in the session
func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
}

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// isAllowedImageType determines if image is among types defined
// in map of allowed images
func isAllowedImageType(mimeType string) bool {
	_, exists := validImageTypes[mimeType]

	return exists
}
