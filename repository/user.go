package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ivamshky/go-crud/common"
	"github.com/ivamshky/go-crud/model"
	"github.com/ivamshky/go-crud/repository/db"
)

type UserRepository interface {
	FindById(ctx context.Context, id int64) (model.User, error)
	ListUsers(ctx context.Context, params model.ListUserParams) ([]model.User, error)
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type UserRepositoryImpl struct {
	q *db.Queries
}

func (u *UserRepositoryImpl) FindById(ctx context.Context, id int64) (model.User, error) {
	slog.Info("In FindById..", "id", id)
	user, err := u.q.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, common.ErrNotFound
		}
		return model.User{}, fmt.Errorf("DB error %w", err)
	}
	return model.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (u *UserRepositoryImpl) ListUsers(ctx context.Context, params model.ListUserParams) ([]model.User, error) {
	paramsSqlc := db.ListUsersParams{}
	if params.Id != nil {
		paramsSqlc.ID = sql.NullInt64{Int64: *params.Id, Valid: true}
	}
	if params.Name != nil {
		paramsSqlc.Name = sql.NullString{
			String: *params.Name,
			Valid:  true,
		}
	}
	if params.Email != nil {
		paramsSqlc.Email = sql.NullString{
			String: *params.Email,
			Valid:  true,
		}
	}
	if params.Limit != nil {
		paramsSqlc.Limit = int32(*params.Limit)
	}

	if params.Offset != nil {
		paramsSqlc.Offset = int32(*params.Offset)
	}

	slog.Info("[REPOSITORY] Searching with Parameters", "params", paramsSqlc)
	usersResult, err := u.q.ListUsers(ctx, paramsSqlc)
	if err != nil {
		return nil, fmt.Errorf("db error when listing %w", err)
	}

	users := make([]model.User, len(usersResult))
	for i, user := range usersResult {
		users[i] = model.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		}
	}

	return users, nil
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {

	createUserParam := db.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
	result, err := u.q.CreateUser(ctx, createUserParam)
	if err != nil {
		return model.User{}, fmt.Errorf("db error while insertion %w", err)
	}
	user.Id, err = result.LastInsertId()
	return user, nil
}

func (u *UserRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	err := u.q.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("db error while deleting %w", err)
	}
	return nil
}

func NewUserRepository(database *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		q: db.New(database),
	}
}
