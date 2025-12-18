package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/config"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/handlers"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/middlewares"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/repositories"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/services"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config, db *gorm.DB) {
	userRepo := repositories.NewUserRepo(db)
	tokenRepo := repositories.NewTokenRepo(db)

	jwtSvc := services.NewJWTService(cfg)
	authSvc := services.NewAuthService(userRepo, tokenRepo, jwtSvc)

	h := handlers.NewAuthHandler(authSvc, userRepo)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	protected := r.Group("/")
	protected.Use(middlewares.AuthRequired(jwtSvc))
	{
		protected.GET("/me", h.Me)
	}
}
