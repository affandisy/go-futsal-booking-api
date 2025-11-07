package handler

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
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

	_ "go-futsal-booking-api/docs"
	// _ "go-futsal-booking-api/internal/dto/request"
	// _ "go-futsal-booking-api/internal/dto/response"
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

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking for a specific schedule (Customer or Admin)
// @Tags Bookings
// @Accept json
// @Produce json
// @Param booking body request.CreateBookingRequest true "Booking creation request"
// @Success 201 {object} docs.SuccessResponse{data=dto.BookingResponse} "Booking successfully created"
// @Failure 400 {object} docs.ErrorResponse "Bad Request or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Schedule Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /bookings [post]
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

	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		logger.Error("Invalid or missing user token")
		return c.JSON(http.StatusUnauthorized, jsonres.Error(
			"UNAUTHORIZED", "Invalid or missing user token", nil,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	newBooking, err := h.bookingService.CreateBooking(
		ctx,
		&request.CreateBookingRequest{
			// UserID:      userID,
			ScheduleID:  req.ScheduleID,
			BookingDate: req.BookingDate,
		},
		userID,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrScheduleNotFound) {
			logger.Error("Schedule not found for booking", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Schedule not found",
				map[string]any{"schedule_id": req.ScheduleID},
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

// GetMyBookings godoc
// @Summary Get bookings by user ID
// @Description Get a list of all bookings for a specific user (e.g., "My Bookings")
// @Tags Bookings
// @Produce json
// @Param user_id query uint true "User ID"
// @Success 200 {object} docs.SuccessResponse{data=[]dto.BookingResponse} "Bookings retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Missing or Invalid User ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "User Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /bookings [get]
func (h *BookingHandler) GetMyBookings(c echo.Context) error {
	// userIdStr := c.Param("user_id")
	userIdStr := c.QueryParam("user_id")
	if userIdStr == "" {
		logger.Error("missing user_id query parameter")
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "missing user_id query parameter", map[string]interface{}{"id": userIdStr},
		))
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid user id", map[string]interface{}{"id": userIdStr},
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

		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"User not found",
				map[string]any{"user_id": userId},
			))
		}

		logger.Error("Failed to retrieve booking", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve booking",
			nil,
		))
	}

	bookingResponses := make([]dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		bookingResponses[i] = dto.ToBookingResponse(booking)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully",
		bookingResponses,
	))
}

// GetBookingDetails godoc
// @Summary Get booking details by ID
// @Description Get details of a specific booking by its ID
// @Tags Bookings
// @Produce json
// @Param id path uint true "Booking ID"
// @Success 200 {object} docs.SuccessResponse{data=dto.BookingResponse} "Booking retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Booking ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Booking Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingDetails(c echo.Context) error {
	bookingIdStr := c.Param("id")

	bookingId, err := strconv.ParseUint(bookingIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid booking id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid booking id", map[string]interface{}{"booking_id": bookingIdStr},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	booking, err := h.bookingService.GetBookingByID(ctx, uint(bookingId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrBookingNotFound) {
			logger.Error("Booking not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Booking not found",
				map[string]any{"booking_id": bookingId},
			))
		}

		logger.Error("Failed to find booking by id", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to find booking by id", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Bookings retrieved successfully", dto.ToBookingResponse(booking),
	))
}
