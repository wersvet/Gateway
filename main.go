package main

import (
	"TEST/internal/config"
	"TEST/internal/middleware"
	"TEST/internal/proxy"
	"TEST/internal/routes"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	authProxy, err := proxy.NewReverseProxy(cfg.AuthServiceURL)
	if err != nil {
		log.Fatalf("failed to create auth proxy: %v", err)
	}
	userProxy, err := proxy.NewReverseProxy(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("failed to create user proxy: %v", err)
	}
	chatProxy, err := proxy.NewReverseProxy(cfg.ChatServiceURL)
	if err != nil {
		log.Fatalf("failed to create chat proxy: %v", err)
	}
	wsProxy, err := proxy.NewWebsocketProxy(cfg.ChatServiceURL)
	if err != nil {
		log.Fatalf("failed to create websocket proxy: %v", err)
	}

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.JWTAuth(cfg.JWTSecret))

	routes.RegisterAuthRoutes(router, authProxy)
	routes.RegisterUserRoutes(router, userProxy)
	routes.RegisterChatRoutes(router, chatProxy, wsProxy)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("API Gateway listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Println("server exiting")
}
