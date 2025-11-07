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
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	_ "go-futsal-booking-api/docs"
	// _ "go-futsal-booking-api/internal/dto/request"
	// _ "go-futsal-booking-api/internal/dto/response"
)

type UserHandler struct {
	userService service.UserService
	timeout     time.Duration
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		timeout:     30 * time.Second,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new customer account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body request.UserRegisterRequest true "User registration details"
// @Success 201 {object} docs.SuccessResponse{data=dto.UserResponse} "User registered successfully"
// @Failure 400 {object} docs.ErrorResponse "Invalid request body or validation error"
// @Failure 400 {object} docs.ErrorResponse "REGISTER_FAILED (e.g., email already exists)"
// @Router /users/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var reqUser request.UserRegisterRequest

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

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

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

	return c.JSON(http.StatusCreated, jsonres.Success(
		"User registered successfully", dto.ToUserResponse(&user),
	))
}

// Login godoc
// @Summary Log in a user
// @Description Log in with email and password to receive a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body request.UserLoginRequest true "User login credentials"
// @Success 200 {object} docs.SuccessResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} docs.ErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} docs.ErrorResponse "LOGIN_FAILED (Invalid credentials or email not verified)"
// @Router /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var reqUser request.UserLoginRequest

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

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

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

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify a user's email account using the provided code from the email link
// @Tags Users
// @Produce json
// @Param code path string true "Email verification code"
// @Success 200 {object} docs.SuccessResponse "Success to Verifying Email"
// @Failure 401 {object} docs.ErrorResponse "Invalid or expired URL"
// @Failure 500 {object} docs.ErrorResponse "Internal server error"
// @Router /users/email-verification/{code} [get]
func (h *UserHandler) VerifyEmail(c echo.Context) error {
	encCode := c.Param("code")

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	err := h.userService.VerifyEmail(ctx, encCode)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return c.JSON(http.StatusUnauthorized, jsonres.Error(
				"INVALID", err.Error(), nil,
			))
		}
		return c.JSON(http.StatusInternalServerError, jsonres.Error(
			"INTERNAL_SERVER_ERROR", err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, jsonres.Success(
		"Success to Verifying Email", nil,
	))
}
