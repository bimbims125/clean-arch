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
}

type CategoryHandler struct {
	Service CategoryService
}

func NewCategoryHandler(r *mux.Router, service CategoryService) {
	handler := &CategoryHandler{Service: service}

	r.HandleFunc("/categories", handler.Fetch).Methods("GET")
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
