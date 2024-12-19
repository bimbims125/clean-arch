package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/gorilla/mux"
)

// ResponseData represent the response data struct
type ResponseData struct {
	Data interface{} `json:"data"`
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// UserService represent the user's usecases

type UserService interface {
	Fetch(ctx context.Context) (result []domain.User, err error)
}

// UserHandler represent the http handler for user
type UserHandler struct {
	Service UserService
}

// NewUserHandler initializes the user HTTP handler
func NewUserHandler(r *mux.Router, service UserService) {
	handler := &UserHandler{Service: service}

	r.HandleFunc("/users", handler.FetchUser).Methods("GET")
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
		json.NewEncoder(w).Encode(ResponseError{Message: err.Error()})
		return
	}

	// Respond with the fetched users in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseData{Data: users})
}
