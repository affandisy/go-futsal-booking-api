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

// GetFieldByID godoc
// @Summary Get field by ID
// @Description Get details of a specific field by its ID
// @Tags Fields
// @Produce json
// @Param id path uint true "Field ID"
// @Success 200 {object} docs.SuccessResponse{data=dto.FieldResponse} "Field retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Field ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Field Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /fields/{id} [get]
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

// GetFieldsByVenue godoc
// @Summary Get fields by venue ID
// @Description Get a list of all fields associated with a specific venue
// @Tags Fields
// @Produce json
// @Param venueId query uint true "Venue ID"
// @Success 200 {object} docs.SuccessResponse{data=[]dto.FieldResponse} "Field retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Missing or Invalid Venue ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Venue Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /fields [get]
func (h *FieldHandler) GetFieldsByVenue(c echo.Context) error {
	// venueId := c.Param("venueId")
	venueId := c.QueryParam("venueId")

	if venueId == "" {
		logger.Error("Missing venueId query parameter")
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Missing 'venueId' query parameter", nil,
		))
	}

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

// CreateField godoc
// @Summary Create a new field (Admin only)
// @Description Create a new field for a specific venue
// @Tags Fields
// @Accept json
// @Produce json
// @Param field body request.CreateFieldRequest true "Field creation request"
// @Success 201 {object} docs.SuccessResponse{data=dto.FieldResponse} "Field successfully created"
// @Failure 400 {object} docs.ErrorResponse "Bad Request or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Venue Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /fields [post]
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

// UpdateField godoc
// @Summary Update a field (Admin only)
// @Description Update details of an existing field
// @Tags Fields
// @Accept json
// @Produce json
// @Param id path uint true "Field ID"
// @Param field body request.UpdateFieldRequest true "Field update request"
// @Success 200 {object} docs.SuccessResponse{data=dto.FieldResponse} "Field successfully updated"
// @Failure 400 {object} docs.ErrorResponse "Bad Request, Invalid ID, or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Field Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /fields/{id} [put]
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

// DeleteField godoc
// @Summary Delete a field (Admin only)
// @Description Delete an existing field by its ID
// @Tags Fields
// @Produce json
// @Param id path uint true "Field ID"
// @Success 200 {object} docs.SuccessResponse "Field deleted successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Field ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Field Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /fields/{id} [delete]
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
