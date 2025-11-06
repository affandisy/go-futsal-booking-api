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

func (h *ScheduleHandler) GetScheduleByID(c echo.Context) error {
	scheduleIdStr := c.Param("id")

	scheduleId, err := strconv.ParseUint(scheduleIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", map[string]interface{}{"id": c.Param("id")},
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

func (h *ScheduleHandler) GetScheduleByField(c echo.Context) error {
	fieldIdStr := c.Param("id")

	fieldId, err := strconv.ParseUint(fieldIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid venue id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid venue id", map[string]interface{}{"id": c.Param("id")},
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

		logger.Error("Failed to create schedule", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to create schedule", nil,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Schedule successfully created", dto.ToScheduleResponse(newSchedule),
	))
}

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
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", err,
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

		logger.Error("Failed to update field", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_ERROR", "Failed to update field", nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule successfully updated", dto.ToScheduleResponse(updatedSchedule),
	))
}

func (h *ScheduleHandler) DeleteSchedule(c echo.Context) error {
	scheduleIdStr := c.Param("id")

	scheduleId, err := strconv.ParseUint(scheduleIdStr, 10, 64)
	if err != nil {
		logger.Error("Invalid field id", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid field id", err,
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

		logger.Error("Failed to delete field", err)
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", "Failed to delete field", err,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Schedule deleted successfully",
		map[string]any{"schedule_id": scheduleId},
	))
}
