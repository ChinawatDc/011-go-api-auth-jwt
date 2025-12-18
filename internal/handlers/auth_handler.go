package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/middlewares"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/repositories"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/services"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/utils"
)

type AuthHandler struct {
	auth  *services.AuthService
	users *repositories.UserRepo
}

func NewAuthHandler(auth *services.AuthService, users *repositories.UserRepo) *AuthHandler {
	return &AuthHandler{auth: auth, users: users}
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	u, err := h.auth.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	access, refresh, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	access, err := h.auth.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	utils.Success(c, gin.H{
		"access_token": access,
		"token_type":   "Bearer",
	})

}

type logoutReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.auth.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uidAny, _ := c.Get(middlewares.CtxUserIDKey)
	uid, _ := uidAny.(uint)

	u, err := h.users.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email})
}
