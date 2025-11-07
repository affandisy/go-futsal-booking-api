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

type ScheduleHandler struct {
	scheduleService service.ScheduleService
	timeout         time.Duration
}

func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
		timeout:         30 * time.Second,
	}
}

// GetScheduleByID godoc
// @Summary Get schedule by ID
// @Description Get details of a specific schedule by its ID
// @Tags Schedules
// @Produce json
// @Param id path uint true "Schedule ID"
// @Success 200 {object} docs.SuccessResponse{data=dto.ScheduleResponse} "Schedule retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Schedule ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Schedule Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedules/{id} [get]
func (h *ScheduleHandler) GetScheduleByID(c echo.Context) error {
	scheduleIdStr := c.Param("id")

	scheduleId, err := strconv.ParseUint(scheduleIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", map[string]interface{}{"id": c.Param("id")},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	schedule, err := h.scheduleService.GetScheduleByID(ctx, uint(scheduleId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("request timeout", map[string]any{"timeout": h.timeout})
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrScheduleNotFound) {
			logger.Error("schedule not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"schedule not found",
				map[string]interface{}{"schedule_id": scheduleId},
			))
		}

		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve schedule",
			nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule retrieved successfully",
		dto.ToScheduleResponse(schedule),
	))
}

// GetScheduleByField godoc
// @Summary Get schedules by field ID
// @Description Get a list of all schedules associated with a specific field
// @Tags Schedules
// @Produce json
// @Param fieldId query uint true "Field ID"
// @Success 200 {object} docs.SuccessResponse{data=[]dto.ScheduleResponse} "Schedule retrieved successfully"
// @Failure 400 {object} docs.ErrorResponse "Missing or Invalid Field ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 404 {object} docs.ErrorResponse "Field Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedules [get]
func (h *ScheduleHandler) GetScheduleByField(c echo.Context) error {
	// fieldIdStr := c.Param("id")

	fieldIdStr := c.QueryParam("fieldId")

	if fieldIdStr == "" {
		logger.Error("invalid fieldId query parameter")
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", nil,
		))
	}

	fieldId, err := strconv.ParseUint(fieldIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", map[string]interface{}{"id": c.Param("id")},
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	schedules, err := h.scheduleService.GetScheduleByField(ctx, uint(fieldId))
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
			logger.Error("field not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"field not found",
				map[string]interface{}{"field_id": fieldId},
			))
		}

		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR",
			"Failed to retrieve schedule",
			nil,
		))
	}

	scheduleResponse := make([]dto.ScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		scheduleResponse[i] = dto.ToScheduleResponse(schedule)
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule retrieved successfully",
		scheduleResponse,
	))
}

// CreateSchedule godoc
// @Summary Create a new schedule (Admin only)
// @Description Create a new weekly schedule for a specific field
// @Tags Schedules
// @Accept json
// @Produce json
// @Param schedule body request.CreateScheduleRequest true "Schedule creation request"
// @Success 201 {object} docs.SuccessResponse{data=dto.ScheduleResponse} "Schedule successfully created"
// @Failure 400 {object} docs.ErrorResponse "Bad Request or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Field Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedules [post]
func (h *ScheduleHandler) CreateSchedule(c echo.Context) error {
	var req request.CreateScheduleRequest

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

	newSchedule, err := h.scheduleService.CreateSchedule(
		ctx,
		&request.CreateScheduleRequest{
			FieldID:   req.FieldID,
			DayOfWeek: req.DayOfWeek,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Price:     req.Price,
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

		if errors.Is(err, domain.ErrFieldNotFound) {
			logger.Error("field not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"field not found",
				map[string]interface{}{"field_id": req.FieldID},
			))
		}

		logger.Error("Failed to create schedule", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to create schedule", nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Schedule successfully created", dto.ToScheduleResponse(newSchedule),
	))
}

// UpdateSchedule godoc
// @Summary Update a schedule (Admin only)
// @Description Update details of an existing schedule
// @Tags Schedules
// @Accept json
// @Produce json
// @Param id path uint true "Schedule ID"
// @Param schedule body request.UpdateScheduleRequest true "Schedule update request"
// @Success 200 {object} docs.SuccessResponse{data=dto.ScheduleResponse} "Schedule successfully updated"
// @Failure 400 {object} docs.ErrorResponse "Bad Request, Invalid ID, or Validation Error"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Schedule Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedules/{id} [put]
func (h *ScheduleHandler) UpdateSchedule(c echo.Context) error {
	scheduleIdStr := c.Param("id")

	var req request.UpdateScheduleRequest
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

	scheduleId, err := strconv.ParseUint(scheduleIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid schedule id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid schedule id", err,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	updatedSchedule, err := h.scheduleService.UpdateSchedule(
		ctx,
		uint(scheduleId),
		req.DayOfWeek,
		req.StartTime,
		req.EndTime,
		req.Price,
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
			logger.Error("schedule not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"schedule not found",
				map[string]interface{}{"schedule_id": scheduleId},
			))
		}

		logger.Error("Failed to update schedule", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to update schedule", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule successfully updated", dto.ToScheduleResponse(updatedSchedule),
	))
}

// DeleteSchedule godoc
// @Summary Delete a schedule (Admin only)
// @Description Delete an existing schedule by its ID
// @Tags Schedules
// @Produce json
// @Param id path uint true "Schedule ID"
// @Success 200 {object} docs.SuccessResponse "Schedule deleted successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid Schedule ID"
// @Failure 401 {object} docs.ErrorResponse "Unauthorized (Missing Token)"
// @Failure 403 {object} docs.ErrorResponse "Forbidden (Not Admin)"
// @Failure 404 {object} docs.ErrorResponse "Schedule Not Found"
// @Failure 500 {object} docs.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedules/{id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c echo.Context) error {
	scheduleIdStr := c.Param("id")

	scheduleId, err := strconv.ParseUint(scheduleIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid schedule id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid schedule id", err,
		))
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	err = h.scheduleService.DeleteSchedule(ctx, uint(scheduleId))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return c.JSON(http.StatusRequestTimeout, jsonres.Error(
				"TIMEOUT",
				"Request timeout",
				nil,
			))
		}

		if errors.Is(err, domain.ErrScheduleNotFound) {
			logger.Error("schedule not found", err)
			return c.JSON(http.StatusNotFound, jsonres.Error(
				"NOT_FOUND",
				"schedule not found",
				map[string]interface{}{"schedule_id": scheduleId},
			))
		}

		logger.Error("Failed to delete schedule", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to delete schedule", err,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule deleted successfully",
		map[string]any{"schedule_id": scheduleId},
	))
}
