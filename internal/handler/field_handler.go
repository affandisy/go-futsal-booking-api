package handler

import (
	"go-futsal-booking-api/internal/dto/request"
	dto "go-futsal-booking-api/internal/dto/response"
	"go-futsal-booking-api/internal/service"
	jsonres "go-futsal-booking-api/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FieldHandler struct {
	fieldService service.FieldService
}

func NewFieldHandler(fieldService service.FieldService) *FieldHandler {
	return &FieldHandler{
		fieldService: fieldService,
	}
}

func (h *FieldHandler) CreateField(c echo.Context) error {
	var req request.CreateFieldRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Failed to fetch request", err,
		))
	}

	newField, err := h.fieldService.CreateField(
		c.Request().Context(),
		req.VenueID,
		req.Name,
		req.FieldType,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Failed to create field", err,
		))
	}

	return c.JSON(http.StatusCreated, jsonres.Success(
		"Field successfully created", dto.ToFieldResponse(&newField),
	))
}
