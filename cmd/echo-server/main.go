package main

import (
	"context"
	"fmt"
	"go-futsal-booking-api/cmd/router"
	"go-futsal-booking-api/internal/handler"
	"go-futsal-booking-api/internal/middleware"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/internal/service"
	"go-futsal-booking-api/pkg/config"
	"go-futsal-booking-api/pkg/database"
	"go-futsal-booking-api/pkg/logger"
	"go-futsal-booking-api/pkg/validator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.App.Environment)
	logger.Info("Starting Futsal Booking API", "version", cfg.App.Version)

	db, err := database.InitPostgres(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}

	logger.Info("Database connected successfully")

	// Init notification from mailjet
	mailjetEmail := repository.NewMailjetRepository(
		repository.MailjetConfig{
			MailjetBaseURL:           cfg.Mailjet.MailjetBaseUrl,
			MailjetBasicAuthUsername: cfg.Mailjet.MailjetBasicAuthUsername,
			MailjetBasicAuthPassword: cfg.Mailjet.MailjetBasicAuthPassword,
			MailjetSenderEmail:       cfg.Mailjet.MailjetSenderEmail,
			MailjetSenderName:        cfg.Mailjet.MailjetSenderName,
		},
	)

	// Init validate
	validate := validator.New()

	// Init repo
	userRepo := repository.NewUserRepository(db)
	fieldRepo := repository.NewFieldRepository(db)
	venueRepo := repository.NewVenueRepository(db)

	// Init service
	userService := service.NewUserService(userRepo, validate, mailjetEmail, cfg.App.AppEmailVerificationKey, cfg.App.AppDeploymentUrl)
	fieldService := service.NewFieldService(fieldRepo, venueRepo)
	venueService := service.NewVenueService(venueRepo)

	// Init handler
	userHandler := handler.NewUserHandler(userService)
	fieldHandler := handler.NewFieldHandler(fieldService)
	venueHandler := handler.NewVenueHandler(venueService)

	// Init echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// HTTP error handler
	e.HTTPErrorHandler = middleware.ErrorHandler

	// Global middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Setup routes
	api := e.Group("/api/v1")
	router.SetupUserRoutes(api, userHandler)
	router.SetupFieldRoutes(api, fieldHandler)
	router.SetupVenueRoutes(api, venueHandler)

	// Goroutine server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Server starting", "address", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped")
}
