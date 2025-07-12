package application

import (
	"errors"

	"github.com/danielrios/product-service-go/internal/core/models"
	"github.com/danielrios/product-service-go/internal/core/ports"
)

// ProductService define a estrutura do nosso serviço de aplicação para produtos.
type ProductService struct {
	repo ports.ProductRepository
}

// NewProductService cria e retorna uma nova instância de ProductService.
func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct lida com a lógica de negócio para criar um novo produto.
func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	validatedProduct, err := models.NewProduct(product.ID, product.Name, product.Price)
	if err != nil {
		return nil, err
	}

	err = s.repo.Add(validatedProduct)
	if err != nil {
		return nil, err
	}
	return validatedProduct, nil
}

// GetProductByID lida com a lógica de negócio para buscar um produto por ID.
func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// GetAllProducts lida com a lógica de negócio para obter todos os produtos.
func (s *ProductService) GetAllProducts() ([]*models.Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct lida com a lógica de negócio para atualizar um produto.
func (s *ProductService) UpdateProduct(id string, product *models.Product) (*models.Product, error) {
	if id != product.ID {
		return nil, errors.New("product ID in path does not match ID in body")
	}

	// Apenas valida os dados, sem criar uma nova instância que zeraria o CreatedAt
	_, err := models.NewProduct(product.ID, product.Name, product.Price)
	if err != nil {
		return nil, err
	}

	err = s.repo.Update(product)
	if err != nil {
		return nil, err
	}
	// Após a atualização, busca e retorna a entidade completa do banco de dados.
	return s.repo.GetByID(id)
}

// DeleteProduct lida com a lógica de negócio para excluir um produto.
func (s *ProductService) DeleteProduct(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
