package repository

import (
	"context"

	"a21hc3NpZ25tZW50/app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUserByID(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	return model.User{}, nil // TODO: replace this
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, nil
		} else {
			return user, err
		}
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdateUserByID(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	var user model.User
	err := r.db.WithContext(ctx).Delete(&user, id).Error
	if err != nil {
		return err
	}
	return nil
}
