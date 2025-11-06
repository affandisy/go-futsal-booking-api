package handler

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/dto/request"
	dto "go-futsal-booking-api/internal/dto/response"
	"go-futsal-booking-api/internal/service"
	"go-futsal-booking-api/pkg/logger"
	jsonres "go-futsal-booking-api/pkg/response"
	"go-futsal-booking-api/pkg/validator"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService service.BookingService
	timeout        time.Duration
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		timeout:        30 * time.Second,
	}
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req request.CreateBookingRequest

	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Failed to fetch request", err,
		))
	}

	if errs := validator.Validate(&req); len(errs) > 0 {
		logger.Error("Failed to validate create request", errs)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation failed", errs,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	newBooking, err := h.bookingService.CreateBooking(
		ctx,
		&request.CreateBookingRequest{
			UserID:      req.UserID,
			ScheduleID:  req.ScheduleID,
			BookingDate: req.BookingDate,
		},
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		logger.Error("Failed to create booking", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to create booking", nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Booking successfully created", dto.ToBookingResponse(newBooking),
	))
}

func (h *BookingHandler) GetMyBookings(c echo.Context) error {
	userIdStr := c.Param("user_id")

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid booking id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid booking id", map[string]interface{}{"id": c.Param("id")},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	bookings, err := h.bookingService.GetMyBookings(ctx, uint(userId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve booking",
			nil,
		))
	}

	bookingResponses := make([]interface{}, len(bookings))
	for i, booking := range bookings {
		bookingResponses[i] = dto.ToBookingResponse(booking)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully",
		bookingResponses,
	))
}

func (h *BookingHandler) GetBookingDetails(c echo.Context) error {
	bookingIdStr := c.Param("id")

	bookingId, err := strconv.ParseUint(bookingIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", map[string]interface{}{"booking_id": bookingId},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	bookings, err := h.bookingService.GetBookingByID(ctx, uint(bookingId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		logger.Error("Failed to find field by venue id", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to find field by venue id", err,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully", bookings,
	))
}
