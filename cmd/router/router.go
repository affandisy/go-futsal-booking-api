package router

import (
	"go-futsal-booking-api/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(api *echo.Group, handler *handler.UserHandler) {
	users := api.Group("/users")

	users.GET("/email-verification/:code", handler.VerifyEmail)
	users.POST("/register", handler.Register)
	users.POST("/login", handler.Login)
}

func SetupFieldRoutes(api *echo.Group, handler *handler.FieldHandler, authRequired echo.MiddlewareFunc, adminOnly echo.MiddlewareFunc) {
	fields := api.Group("/fields")
	fields.GET("", handler.GetFieldsByVenue, authRequired)
	fields.GET("/:id", handler.GetFieldByID, authRequired)

	fields.POST("", handler.CreateField, authRequired, adminOnly)
	fields.PUT("/:id", handler.UpdateField, authRequired, adminOnly)
	fields.DELETE("/:id", handler.DeleteField, authRequired, adminOnly)
}

func SetupVenueRoutes(api *echo.Group, handler *handler.VenueHandler, authRequired echo.MiddlewareFunc, adminOnly echo.MiddlewareFunc) {
	venues := api.Group("/venues")
	venues.GET("", handler.GetAllVenues, authRequired)
	venues.GET("/:id", handler.GetVenueByID, authRequired)

	venues.POST("", handler.CreateVenue, authRequired, adminOnly)
	venues.PUT("/:id", handler.UpdateVenue, authRequired, adminOnly)
	venues.DELETE("/:id", handler.DeleteVenue, authRequired, adminOnly)
}

func SetupScheduleRoutes(api *echo.Group, handler *handler.ScheduleHandler, authRequired echo.MiddlewareFunc, adminOnly echo.MiddlewareFunc) {
	schedules := api.Group("/schedules")
	schedules.GET("", handler.GetScheduleByField, authRequired)
	schedules.GET("/:id", handler.GetScheduleByID, authRequired)

	schedules.POST("", handler.CreateSchedule, authRequired, adminOnly)
	schedules.PUT("/:id", handler.UpdateSchedule, authRequired, adminOnly)
	schedules.DELETE("/:id", handler.DeleteSchedule, authRequired, adminOnly)
}

func SetupBookingRoutes(api *echo.Group, handler *handler.BookingHandler, authRequired echo.MiddlewareFunc, adminOnly echo.MiddlewareFunc) {
	bookings := api.Group("/bookings")
	bookings.GET("/:id", handler.GetBookingDetails, authRequired)
	bookings.GET("", handler.GetMyBookings, authRequired)

	bookings.POST("", handler.CreateBooking, authRequired)
}
