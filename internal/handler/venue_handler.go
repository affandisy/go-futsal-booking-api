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
	for i, venue := range venues {
		venueResponse[i] = dto.ToVenueResponse(&venue)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue retrieved successfully", venueResponse,
	))
}

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

		logger.Error("Failed to update venue", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to update venue", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Venue successfully updated", dto.ToVenueResponse(updateVenue),
	))
}

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
