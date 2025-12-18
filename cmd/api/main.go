package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/config"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/db"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/routes"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal("db connect error: ", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, cfg, database)

	addr := ":" + cfg.AppPort
	log.Println("server running at", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
