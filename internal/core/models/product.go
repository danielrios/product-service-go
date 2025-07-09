package models

import (
	"fmt"
	"time"
)

type Product struct {
	ID        string
	Name      string
	Price     float64
	CreatedAt time.Time
}

func NewProduct(id, name string, price float64) (*Product, error) {
	now := time.Now()
	if id == "" {
		return nil, ErrInvalidProductID
	}

	return &Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: now,
	}, nil
}
func (p Product) String() string {
	return fmt.Sprintf("Product(ID: %s, Name: %s, Price: %.2f, CreatedAt: %s)",
		p.ID, p.Name, p.Price, p.CreatedAt.Format(time.RFC3339))
}
