package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/bimbims125/clean-arch/internal/repository"
	"github.com/bimbims125/clean-arch/utils"
	"github.com/gorilla/mux"
)

// UserService represent the user's usecases

type UserService interface {
	Fetch(ctx context.Context) (result []domain.User, err error)
	Create(ctx context.Context, user domain.User) error
}

// UserHandler represent the http handler for user
type UserHandler struct {
	Service UserService
}

// NewUserHandler initializes the user HTTP handler
func NewUserHandler(r *mux.Router, service UserService) {
	handler := &UserHandler{Service: service}

	r.HandleFunc("/users", handler.FetchUser).Methods("GET")
	r.HandleFunc("/register", handler.Create).Methods("POST")
}

// FetchUser handles HTTP GET /users
func (u *UserHandler) FetchUser(w http.ResponseWriter, r *http.Request) {
	// Create a context from the request
	ctx := r.Context()

	// Fetch the users using the service
	users, err := u.Service.Fetch(ctx)
	if err != nil {
		// Respond with an error if fetching fails
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	// Respond with the fetched users in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: users})
}

func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	hashedPwd, err := repository.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	user.Password = hashedPwd

	err = u.Service.Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.ResponseSuccess{Message: "User created successfully"})
}
