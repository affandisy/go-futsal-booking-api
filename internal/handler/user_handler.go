package handler

import (
	"context"
	"go-futsal-booking-api/internal/dto/request"
	dto "go-futsal-booking-api/internal/dto/response"
	"go-futsal-booking-api/internal/service"
	"go-futsal-booking-api/pkg/logger"
	jsonres "go-futsal-booking-api/pkg/response"
	"go-futsal-booking-api/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var reqUser request.UserRegisterRequest
	var ctx context.Context

	if err := c.Bind(&reqUser); err != nil {
		logger.Error("Invalid request body", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if errs := validator.Validate(&reqUser); len(errs) > 0 {
		logger.Error("Failed to validation user register", errs)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation error", errs,
		))
	}

	if reqUser.RoleID != 1 {
		reqUser.RoleID = 2
	}

	user, err := h.userService.Register(
		ctx,
		reqUser.FullName,
		reqUser.Email,
		reqUser.Password,
		reqUser.Age,
		reqUser.Address,
	)
	if err != nil {
		logger.Error("Failed to register user", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"REGISTER_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"User registered successfully", dto.ToUserResponse(&user),
	))
}

func (h *UserHandler) Login(c echo.Context) error {
	var reqUser request.UserLoginRequest
	var ctx context.Context

	if err := c.Bind(&reqUser); err != nil {
		logger.Error("Failed to bind request", err)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"BAD_REQUEST", "Invalid request body", err.Error(),
		))
	}

	if errs := validator.Validate(&reqUser); len(errs) > 0 {
		logger.Error("Failed to validate user login", errs)
		return c.JSON(http.StatusBadRequest, jsonres.Error(
			"VALIDATION_ERROR", "Validation failed", errs,
		))
	}

	token, user, err := h.userService.Login(ctx, reqUser.Email, reqUser.Password)
	if err != nil {
		logger.Error("Failed to login with user", err)
		return c.JSON(http.StatusUnauthorized, jsonres.Error(
			"LOGIN_FAILED", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Login successful",
		dto.LoginResponse{
			Token: token,
			User:  dto.ToUserResponse(&user),
		},
	))
}
