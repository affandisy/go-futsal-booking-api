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

type VenueHandler struct {
	venueService service.VenueService
	timeout      time.Duration
}

func NewVenueHandler(venueService service.VenueService) *VenueHandler {
	return &VenueHandler{
		venueService: venueService,
		timeout:      30 * time.Second,
	}
}

// GetVenueByID godoc
// @Summary Get venue by ID
// @Description Get details of a specific venue by its ID
// @Tags Venues
// @Produce json
// @Param id path uint true "Venue ID"
// @Success 200 {object} docs.SuccessResponse{data=dto.VenueResponse} "Venue retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Venue ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Venue Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /venues/{id} [get]
func (h *VenueHandler) GetVenueByID(c echo.Context) error {
	venueIdStr := c.Param("id")

	venueId, err := strconv.ParseUint(venueIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", map[string]interface{}{"id": c.Param("id")},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	venue, err := h.venueService.GetVenueByID(ctx, uint(venueId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrVenueNotFound) {
			logger.Error("venue not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Venue not found",
				map[string]interface{}{"venue_id": venueId},
			))
		}

		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve venue",
			nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue retrieved successfully",
		dto.ToVenueResponse(venue),
	))
}

// GetAllVenues godoc
// @Summary Get all venues
// @Description Get a list of all available venues
// @Tags Venues
// @Produce json
// @Success 200 {object} docs.SuccessResponse{data=[]dto.VenueResponse} "Venue retrieved successfully"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /venues [get]
func (h *VenueHandler) GetAllVenues(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	venues, err := h.venueService.GetAllVenues(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		logger.Error("Failed to find all venue", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to find all venue", err,
		))
	}

	venueResponse := make([]dto.VenueResponse, len(venues))
	for i := range venues {
		venueResponse[i] = dto.ToVenueResponse(&venues[i])
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue retrieved successfully", venueResponse,
	))
}

// CreateVenue godoc
// @Summary Create a new venue (Admin only)
// @Description Create a new futsal venue
// @Tags Venues
// @Accept json
// @Produce json
// @Param venue body request.CreateVenueRequest true "Venue creation request"
// @Success 201 {object} docs.SuccessResponse{data=dto.VenueResponse} "Venue successfully created"
// @Failure 400 {object} docs.ErrorResponse "Bad Request or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /venues [post]
func (h *VenueHandler) CreateVenue(c echo.Context) error {
	var req request.CreateVenueRequest

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

	newVenue, err := h.venueService.CreateVenue(
		ctx,
		req.Name,
		req.Address,
		req.City,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		logger.Error("Failed to create venue", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to create venue", nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Venue successfully created", dto.ToVenueResponse(newVenue),
	))
}

// UpdateVenue godoc
// @Summary Update a venue (Admin only)
// @Description Update details of an existing venue
// @Tags Venues
// @Accept json
// @Produce json
// @Param id path uint true "Venue ID"
// @Param venue body request.UpdateVenueRequest true "Venue update request"
// @Success 200 {object} docs.SuccessResponse{data=dto.VenueResponse} "Venue successfully updated"
// @Failure 400 {object} docs.ErrorResponse "Bad Request, Invalid ID, or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Venue Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /venues/{id} [put]
func (h *VenueHandler) UpdateVenue(c echo.Context) error {
	venueIdStr := c.Param("id")

	venueId, err := strconv.ParseUint(venueIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", err,
		))
	}

	var req request.UpdateVenueRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Failed to fetch request", err,
		))
	}

	if errs := validator.Validate(&req); len(errs) > 0 {
		logger.Error("Failed to validate update request", errs)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation failed", errs,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	updateVenue, err := h.venueService.UpdateVenue(
		ctx,
		uint(venueId),
		req.Name,
		req.Address,
		req.City,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrVenueNotFound) {
			logger.Error("venue not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Venue not found",
				map[string]interface{}{"venue_id": venueId},
			))
		}

		logger.Error("Failed to update venue", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to update venue", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue successfully updated", dto.ToVenueResponse(updateVenue),
	))
}

// DeleteVenue godoc
// @Summary Delete a venue (Admin only)
// @Description Delete an existing venue by its ID
// @Tags Venues
// @Produce json
// @Param id path uint true "Venue ID"
// @Success 200 {object} docs.SuccessResponse "Venue deleted successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Venue ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Venue Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /venues/{id} [delete]
func (h *VenueHandler) DeleteVenue(c echo.Context) error {
	venueIdStr := c.Param("id")

	venueId, err := strconv.ParseUint(venueIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", err,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	err = h.venueService.DeleteVenue(ctx, uint(venueId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrVenueNotFound) {
			logger.Error("venue not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Venue not found",
				map[string]interface{}{"venue_id": venueId},
			))
		}

		logger.Error("Failed to delete venue", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to delete venue", err,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue deleted successfully",
		map[string]any{"venue_id": venueId},
	))
}
