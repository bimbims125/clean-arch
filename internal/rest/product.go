package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/bimbims125/clean-arch/utils"
	"github.com/gorilla/mux"
)

type ProductService interface {
	Fetch(ctx context.Context) (result []domain.Product, err error)
}

type ProductHandler struct {
	Service ProductService
}

func NewProductHandler(r *mux.Router, service ProductService) {
	handler := &ProductHandler{Service: service}

	r.HandleFunc("/products", handler.FetchProduct).Methods("GET")
}

func (u *ProductHandler) FetchProduct(w http.ResponseWriter, r *http.Request) {
	// Create a context from the request
	ctx := r.Context()

	// Fetch the products using the service
	products, err := u.Service.Fetch(ctx)
	if err != nil {
		// Respond with an error if fetching fails
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: err.Error()})
		return
	}

	// Respond with the fetched products in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: products})
}
