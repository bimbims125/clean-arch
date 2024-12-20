package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/bimbims125/clean-arch/utils"
	"github.com/gorilla/mux"
)

type ProductService interface {
	Fetch(ctx context.Context) (result []domain.Product, err error)
	FetchPaginated(ctx context.Context, offset, limit int) (total int, products []domain.Product, err error)
	GetByID(ctx context.Context, id int) (result domain.Product, err error)
}

type ProductHandler struct {
	Service ProductService
}

func NewProductHandler(r *mux.Router, service ProductService) {
	handler := &ProductHandler{Service: service}

	r.HandleFunc("/products", handler.FetchPaginatedProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetByID).Methods("GET")
}

func (p *ProductHandler) FetchProduct(w http.ResponseWriter, r *http.Request) {
	// Create a context from the request
	ctx := r.Context()

	// Fetch the products using the service
	products, err := p.Service.Fetch(ctx)
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

func (p *ProductHandler) FetchPaginatedProduct(w http.ResponseWriter, r *http.Request) {
	// Context
	ctx := r.Context()

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	// Default pagination values
	page := 1
	perPage := 10

	// Convert parameters to integers
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	if perPageStr != "" {
		if p, err := strconv.Atoi(perPageStr); err == nil {
			perPage = p
		}
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Fetch data
	total, products, err := p.Service.FetchPaginated(ctx, offset, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate metadata
	totalPages := (total + perPage - 1) / perPage
	metadata := map[string]interface{}{
		"page":        page,
		"per_page":    perPage,
		"sub_total":   len(products),
		"total":       total,
		"total_pages": totalPages,
	}

	// Response
	response := map[string]interface{}{
		"metadata": metadata,
		"products": products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: response})
}

func (p *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Create a context from the request
	ctx := r.Context()

	// Get the product ID from the URL parameters
	id := mux.Vars(r)["id"]

	// Convert ID to integer
	idInt, err := strconv.Atoi(id)

	// Fetch the product using the service
	product, err := p.Service.GetByID(ctx, idInt)
	if err != nil {
		// Handle "not found" error
		if strings.Contains(err.Error(), "not found") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(utils.ResponseError{Message: "product not found"})
			return
		}

		// Handle other errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.ResponseError{Message: "internal server error"})
		return
	}

	// Respond with the fetched product in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.ResponseData{Data: product})
}
