package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ivamshky/go-crud/common"
	"github.com/ivamshky/go-crud/model"
	"github.com/ivamshky/go-crud/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func (u *UserService) GetUserDetails(ctx context.Context, userId int64) (model.User, error) {
	slog.Info("Getting user details", "userId", userId)
	user, err := u.userRepo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return model.User{}, fmt.Errorf("user Not found id: %d", userId)
		}
		return model.User{}, fmt.Errorf("internal Error: %w", err)
	}
	slog.Info("Found user ", "user", user)
	return user, nil
}

func (u *UserService) ListUsers(ctx context.Context, params model.ListUserParams) ([]model.User, error) {
	slog.Info("Finding users", "params", params)
	users, err := u.userRepo.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Internal Error: %w", err)
	}
	slog.Info("Successfully found", "userCount", len(users))
	return users, nil
}

func (u *UserService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	slog.Info("Creating user", "user", user)
	user, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("Internal Error: %w", err)
	}
	slog.Info("Successfully created", "user", user)
	return user, nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId int64) error {
	slog.Info("Deleting user", "userId", userId)
	err := u.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("Internal Error: %w", err)
	}
	slog.Info("Successfully deleted", "user", userId)
	return nil
}
