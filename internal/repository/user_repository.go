package repository

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	UpdateEmailVerification(ctx context.Context, id uint, isVerified bool) error
}

type gormUserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{
		DB: db,
	}
}

func (r *gormUserRepository) Create(ctx context.Context, user *domain.User) error {
	var gormUser gormContract.UserGorm
	gormUser.FromDomain(*user)

	if err := r.DB.WithContext(ctx).Create(&gormUser).Error; err != nil {
		return err
	}

	user.ID = gormUser.ID
	user.CreatedAt = gormUser.CreatedAt

	return nil
}

func (r *gormUserRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var gormUser gormContract.UserGorm

	err := r.DB.WithContext(ctx).Preload("Role").First(&gormUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return gormUser.ToDomain(), nil
}

func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var gormUser gormContract.UserGorm

	err := r.DB.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return gormUser.ToDomain(), nil
}

func (r *gormUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var gormUsers []gormContract.UserGorm

	if err := r.DB.WithContext(ctx).Preload("Role").Find(&gormUsers).Error; err != nil {
		return nil, err
	}

	var users []domain.User
	for _, gu := range gormUsers {
		users = append(users, gu.ToDomain())
	}

	return users, nil
}

func (r *gormUserRepository) Update(ctx context.Context, user *domain.User) error {
	var gormUser gormContract.UserGorm

	if err := r.DB.WithContext(ctx).First(&gormUser, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	gormUser.FullName = user.FullName
	gormUser.Age = user.Age
	gormUser.Address = user.Address

	if err := r.DB.WithContext(ctx).Save(&gormUser).Error; err != nil {
		return err
	}

	user.UpdatedAt = gormUser.UpdatedAt

	return nil
}

func (r *gormUserRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&gormContract.UserGorm{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
}

func (r *gormUserRepository) UpdateEmailVerification(ctx context.Context, id uint, isVerified bool) error {
	result := r.DB.WithContext(ctx).Model(&gormContract.UserGorm{}).Where("id = ?", id).Update("is_verified", isVerified)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or status already updated")
	}

	return nil
}
