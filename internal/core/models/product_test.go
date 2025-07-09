package models_test

import (
	"github.com/danielrios/product-service-go/internal/core/models"
	"strings"
	"testing"
)

func TestNewProduct(t *testing.T) {
	t.Run("Valid Product Creation", func(t *testing.T) {
		id := "1"
		name := "Product 1"
		price := 1234.50
		product, err := models.NewProduct(id, name, price)

		if err != nil {
			t.Errorf("Expected nil, got %s", err)
		}

		if product.ID != id {
			t.Errorf("Expected ID %s, got %s", id, product.ID)
		}

		if product.Name != name {
			t.Errorf("Expected Name %s, got %s", name, product.Name)
		}

		if product.Price != price {
			t.Errorf("Expected Price %.2f, got %.2f", price, product.Price)
		}
	})

	t.Run("Invalid Product", func(t *testing.T) {
		name := ""
		price := 0.0
		product, err := models.NewProduct("", name, price)

		if err == nil {
			t.Error("Expected an error for invalid product creation, got nil")
		}

		if product != nil {
			t.Errorf("Expected product to be nil, got %+v", product)
		}

		if err.Error() != models.ErrInvalidProductID.Error() {
			t.Errorf("Expected error message 'ID cannot be empty', got '%s'", err.Error())
		}
	})

	t.Run("Product String Representation", func(t *testing.T) {
		id := "1"
		name := "Product 1"
		price := 1234.50
		product, _ := models.NewProduct(id, name, price)

		expectedString := "Product(ID: 1, Name: Product 1, Price: 1234.50, CreatedAt: "
		if !strings.Contains(product.String(), expectedString) {
			t.Errorf("Expected string to contain '%s', got '%s'", expectedString, product.String())
		}
	})

}
