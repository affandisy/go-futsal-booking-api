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

	_ "go-futsal-booking-api/docs"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Futsal Booking API
// @version 1.0
// @description This is a REST API for a Futsal Booking application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	scheduleRepo := repository.NewScheduleRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Init service
	userService := service.NewUserService(userRepo, validate, mailjetEmail, cfg.App.AppEmailVerificationKey, cfg.App.AppDeploymentUrl)
	fieldService := service.NewFieldService(fieldRepo, venueRepo, scheduleRepo)
	venueService := service.NewVenueService(venueRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, fieldRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, scheduleRepo, userRepo)

	// Init handler
	userHandler := handler.NewUserHandler(userService)
	fieldHandler := handler.NewFieldHandler(fieldService)
	venueHandler := handler.NewVenueHandler(venueService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	bookingHandler := handler.NewBookingHandler(bookingService)

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

	// Auth middleware
	authRequired := middleware.AuthMiddleware()
	adminOnly := middleware.AdminOnly()

	// Swagger Documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Setup routes
	api := e.Group("/api/v1")
	router.SetupUserRoutes(api, userHandler)
	router.SetupFieldRoutes(api, fieldHandler, authRequired, adminOnly)
	router.SetupVenueRoutes(api, venueHandler, authRequired, adminOnly)
	router.SetupScheduleRoutes(api, scheduleHandler, authRequired, adminOnly)
	router.SetupBookingRoutes(api, bookingHandler, authRequired, adminOnly)

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
