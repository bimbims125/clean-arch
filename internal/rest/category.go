package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/bimbims125/clean-arch/utils"
	"github.com/gorilla/mux"
)

type CategoryService interface {
	Fetch(ctx context.Context) (result []domain.Category, err error)
	GetByID(ctx context.Context, id string) (result domain.Category, err error)
	Create(ctx context.Context, category domain.Category) (err error)
}

type CategoryHandler struct {
	Service CategoryService
}

func NewCategoryHandler(r *mux.Router, service CategoryService) {
	handler := &CategoryHandler{Service: service}

	r.HandleFunc("/categories", handler.Fetch).Methods("GET")
	r.HandleFunc("/categories/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/categories", handler.Create).Methods("POST")
}

func (c *CategoryHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	// Create a context from the request
	ctx := r.Context()

	// Fetch the categories using the service
	categories, err := c.Service.Fetch(ctx)
	if err != nil {
		// Respond with an error if fetching fails
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	// Respond with the fetched categories in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: categories})
}

func (c *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Create a context from a request
	ctx := r.Context()

	// Get the category ID from request
	id := mux.Vars(r)["id"]

	// Get the category by ID using the servuce'
	category, err := c.Service.GetByID(ctx, id)
	if err != nil {
		// Respond with an error if fetching fails
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: category})
}

func (c *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decode json request
	var category domain.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(utils.ResponseError{Message: "Invalid request Payload"})
		return
	}

	err = c.Service.Create(r.Context(), category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.ResponseSuccess{Message: "Category created successfully"})
}
