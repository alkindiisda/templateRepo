package service

import (
	"context"
	"errors"
	"time"

	"a21hc3NpZ25tZW50/app/model"
	"a21hc3NpZ25tZW50/app/repository"
)

type UserService interface {
	Login(ctx context.Context, user *model.User) (id int, err error)
	Register(ctx context.Context, user *model.User) (model.User, error)

	GetByID(ctx context.Context, id int) (*model.User, error)
	UpdateByID(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Login(ctx context.Context, user *model.User) (id int, err error) {
	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return 0, errors.New("user not found")
	}

	if user.Password != dbUser.Password {
		return 0, errors.New("wrong email or password")
	}

	return dbUser.ID, nil
}

func (s *userService) Register(ctx context.Context, user *model.User) (model.User, error) {
	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	user.CreatedAt = time.Now()

	newUser, err := s.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepository.DeleteUser(ctx, id)
}

func (s *userService) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) UpdateByID(ctx context.Context, user *model.User) error {
	return s.userRepository.UpdateUserByID(ctx, user)
}
