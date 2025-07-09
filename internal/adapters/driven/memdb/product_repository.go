package memdb

import (
	"github.com/danielrios/product-service-go/internal/core/models"
	"github.com/danielrios/product-service-go/internal/core/ports"
	"sync"
)

// InMemoryProductRepository é um Adaptador de Saída (Driven Adapter) que implementa a porta ports ProductRepository definida no Core.
type InMemoryProductRepository struct {
	products map[string]*models.Product
	mu       sync.RWMutex
}

// NewInMemoryProductRepository cria uma nova instância do repositório de produtos em memória.
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*models.Product),
	}
}

var _ ports.ProductRepository = (*InMemoryProductRepository)(nil)

// Add adiciona um novo produto ao repositório em memória.
func (r *InMemoryProductRepository) Add(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; ok {
		return models.ErrProductAlreadyExists
	}
	r.products[product.ID] = product
	return nil
}

// GetByID busca um produto pelo seu ID no repositório em memória.
func (r *InMemoryProductRepository) GetByID(id string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, models.ErrProductNotFound
	}
	return product, nil
}

// GetAll retorna todos os produtos armazenados no repositório em memória.
func (r *InMemoryProductRepository) GetAll() ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allProducts := make([]*models.Product, 0)
	for _, p := range r.products {
		allProducts = append(allProducts, p)
	}
	return allProducts, nil
}

// Update atualiza um produto existente no repositório em memória.
func (r *InMemoryProductRepository) Update(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; !ok {
		return models.ErrProductNotFound
	}
	r.products[product.ID] = product
	return nil
}

// Delete remove um produto pelo seu ID do repositório em memória.
func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return models.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}
