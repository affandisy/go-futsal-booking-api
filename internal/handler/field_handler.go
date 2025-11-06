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
)

type FieldHandler struct {
	fieldService service.FieldService
	timeout      time.Duration
}

func NewFieldHandler(fieldService service.FieldService) *FieldHandler {
	return &FieldHandler{
		fieldService: fieldService,
		timeout:      30 * time.Second,
	}
}

func (h *FieldHandler) GetFieldByID(c echo.Context) error {
	fieldIdStr := c.Param("id")

	fieldId, err := strconv.ParseUint(fieldIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", map[string]interface{}{"id": c.Param("id")},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	field, err := h.fieldService.GetFieldByID(ctx, uint(fieldId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrFieldNotFound) {
			logger.Error("Field not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Field not found",
				map[string]any{"field_id": fieldId},
			))
		}

		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve field",
			nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Field retrieved successfully",
		dto.ToFieldResponse(field),
	))
}

func (h *FieldHandler) GetFieldsByVenue(c echo.Context) error {
	venueId := c.Param("venueId")

	venueIdUint, err := strconv.ParseUint(venueId, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", map[string]interface{}{"venue_id": venueIdUint},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	fields, err := h.fieldService.GetFieldsByVenue(ctx, uint(venueIdUint))
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
			logger.Error("Field not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Venue not found",
				map[string]interface{}{"venue_id": venueIdUint},
			))
		}

		logger.Error("Failed to find field by venue id", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to find field by venue id", err,
		))
	}

	fieldResponse := make([]dto.FieldResponse, len(fields))
	for i, field := range fields {
		fieldResponse[i] = dto.ToFieldResponse(field)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Field retrieved successfully", fieldResponse,
	))
}

func (h *FieldHandler) CreateField(c echo.Context) error {
	var req request.CreateFieldRequest

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

	newField, err := h.fieldService.CreateField(
		ctx,
		&request.CreateFieldRequest{
			VenueID:   req.VenueID,
			Name:      req.Name,
			FieldType: req.FieldType,
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

		if errors.Is(err, domain.ErrVenueNotFound) {
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Venue not found",
				map[string]interface{}{"venue_id": req.VenueID},
			))
		}

		logger.Error("Failed to create field", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to create field", nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Field successfully created", dto.ToFieldResponse(newField),
	))
}

func (h *FieldHandler) UpdateField(c echo.Context) error {
	fieldId := c.Param("id")

	var req request.UpdateFieldRequest
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

	fieldIdInt, err := strconv.ParseUint(fieldId, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", err,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	updatedField, err := h.fieldService.UpdateField(
		ctx,
		uint(fieldIdInt),
		req.Name,
		req.FieldType,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrFieldNotFound) {
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Field not found",
				map[string]interface{}{"field_id": fieldIdInt},
			))
		}

		logger.Error("Failed to update field", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to update field", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Field successfully updated", dto.ToFieldResponse(updatedField),
	))
}

func (h *FieldHandler) DeleteField(c echo.Context) error {
	fieldId := c.Param("id")

	fieldIdInt, err := strconv.ParseUint(fieldId, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", err,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	err = h.fieldService.DeleteField(
		ctx,
		uint(fieldIdInt),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrFieldNotFound) {
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"Field not found",
				map[string]interface{}{"field_id": fieldIdInt},
			))
		}

		logger.Error("Failed to delete field", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to delete field", err,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Field deleted successfully",
		map[string]any{"field_id": fieldId},
	))
}
