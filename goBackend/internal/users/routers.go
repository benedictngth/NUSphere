package users

//contains the handlers for the user routes and routers

import (
	// "errors"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Users(router *gin.RouterGroup, authService AuthService) {
	router.GET("", GetUsersHandler(authService))
	router.POST("/register", RegisterHandler(authService))
	router.POST("/login", LoginHandler(authService))
	router.GET("/cookie", GetCookieHandler(authService))
	router.GET("/deleteCookie", DeleteCookieHandler(authService))
}
func GetCookieHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("Authorisation")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cookie not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"cookie": cookie})
	}
}

func DeleteCookieHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorisation", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"cookie": "deleted"})
	}
}

func Profile(router *gin.RouterGroup) {
	router.GET("/profile", func(c *gin.Context) {
		userID, userErr := c.Get("user_id")
		if !userErr {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userID})
	})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// handler for registering a new user in the system
func RegisterHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user, err := authService.Register(context.Background(), req.Username, req.Password)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": user.PublicID, "username": user.Username})

	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handler for logging in a user calls authService.Login
func LoginHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		token, err := authService.Login(context.Background(), req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		log.Print("generated token string: ", token)
		c.SetCookie("Authorisation", fmt.Sprintf("Bearer %v", token), 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, "logged in")
	}
}

func GetUsersHandler(authService AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := authService.GetUsers(context.Background())
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get users"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}