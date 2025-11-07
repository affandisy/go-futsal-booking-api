package service_test

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository/mock"
	"go-futsal-booking-api/internal/service"
	"go-futsal-booking-api/pkg/logger"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize logger for testing
	logger.Init("test")
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockNotifRepo := mock.NewMockNotificationRepository(ctrl)
	validate := validator.New()

	userService := service.NewUserService(
		mockUserRepo,
		validate,
		mockNotifRepo,
		"test-encryption-key-32-characters",
		"http://localhost:8080",
	)

	t.Run("Success - Register new user", func(t *testing.T) {
		ctx := context.Background()
		fullName := "John Doe"
		email := "john.doe@example.com"
		password := "password123"
		age := 25
		address := "123 Main St"

		// Mock: email should not exist
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(domain.User{}, errors.New("user not found"))

		// Mock: create user
		mockUserRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *domain.User) error {
				user.ID = 1
				user.CreatedAt = time.Now()
				return nil
			})

		// Mock: send email (can fail without affecting the flow)
		mockNotifRepo.EXPECT().
			SendEmail(gomock.Any(), email, gomock.Any(), gomock.Any()).
			Return(nil)

		result, err := userService.Register(ctx, fullName, email, password, age, address)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, fullName, result.FullName)
		assert.Equal(t, email, result.Email)
		assert.Equal(t, "", result.Password) // Password should be cleared
		assert.False(t, result.IsVerified)
	})

	t.Run("Fail - Email already exists", func(t *testing.T) {
		ctx := context.Background()
		email := "existing@example.com"

		// Mock: email already exists
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(domain.User{ID: 1, Email: email}, nil)

		result, err := userService.Register(ctx, "John Doe", email, "password123", 25, "Address")

		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Invalid email format", func(t *testing.T) {
		ctx := context.Background()
		invalidEmail := "invalid-email"

		result, err := userService.Register(ctx, "John Doe", invalidEmail, "password123", 25, "Address")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Password too short", func(t *testing.T) {
		ctx := context.Background()
		email := "test@example.com"
		shortPassword := "123"

		result, err := userService.Register(ctx, "John Doe", email, shortPassword, 25, "Address")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be at least 6 characters")
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Age too young", func(t *testing.T) {
		ctx := context.Background()
		email := "test@example.com"

		result, err := userService.Register(ctx, "John Doe", email, "password123", 10, "Address")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "age must be at leats 15 years old")
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Database error on create", func(t *testing.T) {
		ctx := context.Background()
		email := "test@example.com"

		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(domain.User{}, errors.New("user not found"))

		mockUserRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("database error"))

		result, err := userService.Register(ctx, "John Doe", email, "password123", 25, "Address")

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, uint(0), result.ID)
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockNotifRepo := mock.NewMockNotificationRepository(ctrl)
	validate := validator.New()

	userService := service.NewUserService(
		mockUserRepo,
		validate,
		mockNotifRepo,
		"test-encryption-key-32-characters",
		"http://localhost:8080",
	)

	t.Run("Success - Login with valid credentials", func(t *testing.T) {
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "password123"

		// Use actual bcrypt for testing with cost 10 for speed
		hashedPassword := "$2a$10$RZRAkKRSKe/DR8AaCo8N6e0pJW.eDwsOUMbHkrDoa1OWAkTQ9Y4Oy"

		user := domain.User{
			ID:         1,
			FullName:   "John Doe",
			Email:      email,
			Password:   hashedPassword,
			IsVerified: true,
			Role: domain.Role{
				ID:       2,
				RoleName: "customer",
			},
		}

		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		token, result, err := userService.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, "", result.Password) // Password should be cleared
	})

	t.Run("Fail - User not found", func(t *testing.T) {
		ctx := context.Background()
		email := "nonexistent@example.com"
		password := "password123"

		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(domain.User{}, errors.New("user not found"))

		token, result, err := userService.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Incorrect password", func(t *testing.T) {
		ctx := context.Background()
		email := "john.doe@example.com"
		wrongPassword := "wrongpassword"

		hashedPassword := "$2a$10$RZRAkKRSKe/DR8AaCo8N6e0pJW.eDwsOUMbHkrDoa1OWAkTQ9Y4Oy"

		user := domain.User{
			ID:         1,
			Email:      email,
			Password:   hashedPassword,
			IsVerified: true,
			Role: domain.Role{
				ID:       2,
				RoleName: "customer",
			},
		}

		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		token, result, err := userService.Login(ctx, email, wrongPassword)

		assert.Error(t, err)
		assert.Equal(t, "incorrect password", err.Error())
		assert.Empty(t, token)
		assert.Equal(t, uint(0), result.ID)
	})

	t.Run("Fail - Email not verified", func(t *testing.T) {
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "password123"

		hashedPassword := "$2a$10$RZRAkKRSKe/DR8AaCo8N6e0pJW.eDwsOUMbHkrDoa1OWAkTQ9Y4Oy"

		user := domain.User{
			ID:         1,
			Email:      email,
			Password:   hashedPassword,
			IsVerified: false, // Not verified
			Role: domain.Role{
				ID:       2,
				RoleName: "customer",
			},
		}

		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		token, result, err := userService.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Equal(t, "email address has not been verified", err.Error())
		assert.Empty(t, token)
		assert.Equal(t, uint(0), result.ID)
	})
}
