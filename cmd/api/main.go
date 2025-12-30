package main

import (
	"log"

	"komando/internal/config"
	"komando/internal/db"
	"komando/internal/middlewares"
	"komando/internal/modules/auth"
	"komando/internal/shared/response"

	docs "komando/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Komando API
// @version 1.0
// @description Backend Komando (Check-in, Materials, Quiz, Reminders)
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DBDSN == "" {
		log.Fatal("DB_DSN is empty")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is empty")
	}

	conn, err := db.NewPostgres(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// swagger metadata
	docs.SwaggerInfo.BasePath = "/api/v1"

	// swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// health
	r.GET("/health", func(c *gin.Context) { response.OK(c, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")

	// DI: auth
	authRepo := auth.NewRepo(conn)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiresMin)
	authH := auth.NewHandler(authSvc)

	api.POST("/auth/login", authH.Login)

	// protected
	protected := api.Group("/")
	protected.Use(middlewares.AuthJWT(cfg.JWTSecret))
	protected.GET("/auth/me", authH.Me)

	log.Printf("Komando API running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
