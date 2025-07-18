package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/danielrios/product-service-go/internal/application"
	"github.com/danielrios/product-service-go/internal/core/models"
	"github.com/go-chi/chi/v5"
)

// ProductHandler define a estrutura do nosso Adaptador de Entrada HTTP.
type ProductHandler struct {
	service *application.ProductService
}

// NewProductHandler cria e retorna uma nova instância de ProductHandler.
func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// writeJSONResponse é um helper para enviar respostas JSON padronizadas.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Erro ao codificar resposta JSON: %v", err)
		}
	}
}

// writeErrorResponse é um helper para enviar respostas de erro padronizadas.
func writeErrorResponse(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	if errors.Is(err, models.ErrProductNotFound) {
		statusCode = http.StatusNotFound
		message = err.Error()
	} else if errors.Is(err, models.ErrProductAlreadyExists) {
		statusCode = http.StatusConflict
		message = err.Error()
	} else if errors.Is(err, models.ErrInvalidProductID) {
		statusCode = http.StatusBadRequest
		message = err.Error()
	} else {
		log.Printf("Erro interno não mapeado no handler: %v", err)
	}

	writeJSONResponse(w, statusCode, map[string]string{"error": message})
}

// CreateProductHandler lida com a requisição POST /products.
func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeErrorResponse(w, errors.New("invalid request body"))
		return
	}

	createdProduct, err := h.service.CreateProduct(&product)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdProduct)
}

// GetProductByIDHandler lida com a requisição GET /products/{id}.
func (h *ProductHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.service.GetProductByID(id)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, product)
}

// GetAllProductsHandler lida com a requisição GET /products (listagem)
func (h *ProductHandler) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, products)
}

// UpdateProductHandler lida com a requisição PUT /products/{id}
func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeErrorResponse(w, errors.New("invalid request body"))
		return
	}

	updatedProduct, err := h.service.UpdateProduct(id, &product)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedProduct)
}

// DeleteProductHandler lida com a requisição DELETE /products/{id}
func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteProduct(id)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
