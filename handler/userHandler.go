package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ivamshky/go-crud/common"
	"github.com/ivamshky/go-crud/model"
	"github.com/ivamshky/go-crud/service"
)

var (
	DEFAULT_LIMIT  int = 10
	DEFAULT_OFFSET int = 0
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s: s}
}

func (h *UserHandler) HandleGetUserById(writer http.ResponseWriter, request *http.Request) {
	rawUserId := request.PathValue("userId")
	var userId int64
	_, err := fmt.Sscanf(rawUserId, "%d", &userId)
	if err != nil {
		slog.Error("Invalid userID received: ", "userId", rawUserId)
		http.Error(writer, "Invalid userId format. expecting /user/{id}", http.StatusBadRequest)
		return
	}

	user, err := h.s.GetUserDetails(request.Context(), userId)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			http.Error(writer, fmt.Sprintf("User with Id: %d Not found.", userId), http.StatusNotFound)
			return
		}
		slog.Error("Internal server error: ", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(user); err != nil {
		slog.Error("Error encoding success response: %v", err)
	}
}

func (h *UserHandler) HandleCreateUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var userbody model.User
	err := decoder.Decode(&userbody)
	if err != nil {
		slog.Error("Error Unmarshaling: ", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	createdUser, err := h.s.CreateUser(request.Context(), userbody)
	if err != nil {
		slog.Error("Error Creating user: ", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(createdUser); err != nil {
		slog.Error("Error encoding success response: %v", err)
	}
}

func (h *UserHandler) HandleListUsers(writer http.ResponseWriter, request *http.Request) {
	rawQueryParams := request.URL.Query()
	slog.Info("Raw Query Params", "queryParams", rawQueryParams)
	queryParams := model.ListUserParams{}
	if rawQueryParams.Get("id") != "" {
		if v, err := strconv.ParseInt(rawQueryParams.Get("id"), 10, 64); err == nil {
			queryParams.Id = &v
		} else {
			slog.Error("Error parsing Id: ", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	if name := rawQueryParams.Get("name"); name != "" {
		queryParams.Name = &name
	}

	if email := rawQueryParams.Get("email"); email != "" {
		queryParams.Email = &email
	}

	if limit := rawQueryParams.Get("limit"); limit != "" {
		if limit, err := strconv.Atoi(rawQueryParams.Get("limit")); err == nil {
			queryParams.Limit = &limit
		} else {
			slog.Error("Error parsing limit: ", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
	} else {
		queryParams.Limit = &DEFAULT_LIMIT
	}

	if offset := rawQueryParams.Get("offset"); offset != "" {
		if offset, err := strconv.Atoi(rawQueryParams.Get("offset")); err == nil {
			queryParams.Offset = &offset
		} else {
			slog.Error("Error parsing offset: ", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
	} else {
		queryParams.Offset = &DEFAULT_OFFSET
	}

	listUsers, err := h.s.ListUsers(request.Context(), queryParams)
	if err != nil {
		slog.Error("Error fetching the list: ", err)
		http.Error(writer, "Internal Server error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(listUsers); err != nil {
		slog.Error("Error encoding success response: %v", err)
	}
}
