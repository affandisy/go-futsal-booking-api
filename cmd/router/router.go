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

func SetupFieldRoutes(api *echo.Group, handler *handler.FieldHandler) {
	fields := api.Group("/fields")
	fields.GET("", handler.GetFieldsByVenue)
	fields.GET("/:id", handler.GetFieldByID)

	fields.POST("", handler.CreateField)
	fields.PUT("/:id", handler.UpdateField)
	fields.DELETE("/:id", handler.DeleteField)
}

func SetupVenueRoutes(api *echo.Group, handler *handler.VenueHandler) {
	venues := api.Group("/venues")
	venues.GET("", handler.GetAllVenues)
	venues.GET("/:id", handler.GetVenueByID)

	venues.POST("", handler.CreateVenue)
	venues.PUT("/:id", handler.UpdateVenue)
	venues.DELETE("/:id", handler.DeleteVenue)
}

func SetupScheduleRoutes(api *echo.Group, handler *handler.ScheduleHandler) {
	schedules := api.Group("/schedules")
	schedules.GET("", handler.GetScheduleByField)
	schedules.GET("/:id", handler.GetScheduleByID)

	schedules.POST("", handler.CreateSchedule)
	schedules.PUT("/:id", handler.UpdateSchedule)
	schedules.DELETE("/:id", handler.DeleteSchedule)
}

func SetupBookingRoutes(api *echo.Group, handler *handler.BookingHandler) {
	bookings := api.Group("/bookings")
	bookings.GET("/:id", handler.GetMyBookings)
	bookings.GET("", handler.GetBookingDetails)

	bookings.POST("", handler.CreateBooking)
}
