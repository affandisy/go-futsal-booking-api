package router

import (
	"go-futsal-booking-api/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(api *echo.Group, handler *handler.UserHandler) {
	users := api.Group("/users")
	users.POST("/register", handler.Register)
	users.POST("/login", handler.Login)
}
