package ports

import "github.com/danielrios/product-service-go/internal/core/models"

// ProductRepository define a porta (interface) para operações de persistência de produtos.
// Esta interface é agnóstica a qualquer tecnologia de banco de dados ou forma de armazenamento.
// Ela representa o contrato que o domínio espera de qualquer adaptador de persistência.
type ProductRepository interface {
	GetAll() ([]*models.Product, error)
	GetByID(id string) (*models.Product, error)
	Add(product *models.Product) error
	Update(product *models.Product) error
	Delete(id string) error
}
