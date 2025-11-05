package service

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
	"go-futsal-booking-api/pkg/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Register(ctx context.Context, fullName, email, password string, age int, address string) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	validate *validator.Validate
}

func NewUserService(userRepo repository.UserRepository, validate *validator.Validate) UserService {
	return &userService{
		userRepo: userRepo,
		validate: validate,
	}
}

func (s *userService) Register(ctx context.Context, fullName, email, password string, age int, address string) (domain.User, error) {
	if err := s.validate.Var(email, "required, email"); err != nil {
		logger.Error("Invalid email format", err)
		return domain.User{}, errors.New("invalid email format")
	}

	if err := s.validate.Var(password, "required, min=6"); err != nil {
		logger.Error("Invalid user password", err)
		return domain.User{}, errors.New("password must be at least 6 characters")
	}

	if err := s.validate.Var(age, "required, min=15"); err != nil {
		logger.Error("Invalid user age", err)
		return domain.User{}, errors.New("age must be at leats 15 years old")
	}

	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		logger.Error("Email already exists", err)
		return domain.User{}, errors.New("email already exists")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", err)
		return domain.User{}, errors.New("failed to hash password")
	}

	newUser := domain.User{
		FullName:   fullName,
		Email:      email,
		Password:   string(passwordHash),
		Age:        age,
		Address:    address,
		IsVerified: false,
		Role: domain.Role{
			ID: 2,
		},
	}

	if err := s.userRepo.Create(ctx, &newUser); err != nil {
		logger.Error("Failed to create new user")
		return domain.User{}, err
	}

	newUser.Password = ""
	return newUser, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error("Invalid user credentials", err)
		return "", domain.User{}, err
	}

	ok := utils.CheckPassword(password, user.Password)
	if !ok {
		logger.Error("User password incorrect", err)
		return "", domain.User{}, errors.New("incorrect password")
	}

	userIdStr := strconv.FormatUint(uint64(user.ID), 10)
	token, err := utils.GenerateJWT(userIdStr, user.Role.RoleName)
	if err != nil {
		logger.Error("Failed to generated token", err)
		return "", domain.User{}, errors.New("failed to generate token")
	}

	user.Password = ""
	return token, user, nil
}
