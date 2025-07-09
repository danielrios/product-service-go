package memdb_test

import (
	"errors"
	"testing"

	"github.com/danielrios/product-service-go/internal/adapters/driven/memdb"
	"github.com/danielrios/product-service-go/internal/core/models"
)

func TestInMemoryProductRepository_Add(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, err := models.NewProduct("1", "Test Product", 100.0)
		if err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}

		err = repo.Add(product)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		savedProduct, err := repo.GetByID("1")
		if err != nil {
			t.Errorf("Failed to get added product: %v", err)
		}
		if savedProduct.ID != product.ID || savedProduct.Name != product.Name || savedProduct.Price != product.Price {
			t.Errorf("Saved product doesn't match original. Got %v, want %v", savedProduct, product)
		}
	})

	t.Run("Product Already Exists", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Test Product", 100.0)

		_ = repo.Add(product)

		err := repo.Add(product)

		if !errors.Is(err, models.ErrProductAlreadyExists) {
			t.Errorf("Expected ErrProductAlreadyExists, got %v", err)
		}
	})
}

func TestInMemoryProductRepository_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Test Product", 100.0)
		_ = repo.Add(product)

		retrievedProduct, err := repo.GetByID("1")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if retrievedProduct.ID != product.ID || retrievedProduct.Name != product.Name || retrievedProduct.Price != product.Price {
			t.Errorf("Retrieved product doesn't match original. Got %v, want %v", retrievedProduct, product)
		}
	})

	t.Run("Product Not Found", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()

		_, err := repo.GetByID("nonexistent")

		if !errors.Is(err, models.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got %v", err)
		}
	})
}

func TestInMemoryProductRepository_GetAll(t *testing.T) {
	t.Run("Success With Products", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product1, _ := models.NewProduct("1", "Product 1", 100.0)
		product2, _ := models.NewProduct("2", "Product 2", 200.0)
		_ = repo.Add(product1)
		_ = repo.Add(product2)

		products, err := repo.GetAll()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}

		foundProduct1 := false
		foundProduct2 := false
		for _, p := range products {
			if p.ID == "1" {
				foundProduct1 = true
			}
			if p.ID == "2" {
				foundProduct2 = true
			}
		}
		if !foundProduct1 || !foundProduct2 {
			t.Errorf("Not all products were returned. Found product1: %v, Found product2: %v",
				foundProduct1, foundProduct2)
		}
	})

	t.Run("Success With Empty Repository", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()

		products, err := repo.GetAll()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if products == nil {
			t.Error("Expected empty slice, got nil")
		}
		if len(products) != 0 {
			t.Errorf("Expected 0 products, got %d", len(products))
		}
	})
}

func TestInMemoryProductRepository_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Original Product", 100.0)
		_ = repo.Add(product)

		updatedProduct, _ := models.NewProduct("1", "Updated Product", 150.0)

		err := repo.Update(updatedProduct)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		retrievedProduct, _ := repo.GetByID("1")
		if retrievedProduct.Name != "Updated Product" || retrievedProduct.Price != 150.0 {
			t.Errorf("Product was not updated correctly. Got %v", retrievedProduct)
		}
	})

	t.Run("Product Not Found", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Test Product", 100.0)

		err := repo.Update(product)

		if !errors.Is(err, models.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got %v", err)
		}
	})
}

func TestInMemoryProductRepository_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Test Product", 100.0)
		_ = repo.Add(product)

		err := repo.Delete("1")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		_, err = repo.GetByID("1")
		if !errors.Is(err, models.ErrProductNotFound) {
			t.Errorf("Expected product to be deleted, but it still exists")
		}
	})

	t.Run("Product Not Found", func(t *testing.T) {
		repo := memdb.NewInMemoryProductRepository()

		err := repo.Delete("nonexistent")

		if !errors.Is(err, models.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got %v", err)
		}
	})
}

func TestInMemoryProductRepository_Concurrency(t *testing.T) {
	t.Run("Concurrent Operations", func(t *testing.T) {

		repo := memdb.NewInMemoryProductRepository()
		product, _ := models.NewProduct("1", "Test Product", 100.0)
		_ = repo.Add(product)

		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(index int) {
				if index%2 == 0 {
					_, _ = repo.GetByID("1")
					_, _ = repo.GetAll()
				} else {
					updatedProduct, _ := models.NewProduct("1", "Updated Product", float64(100+index))
					_ = repo.Update(updatedProduct)
				}
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		_, err := repo.GetByID("1")
		if err != nil {
			t.Errorf("Expected product to still exist after concurrent operations")
		}
	})
}
