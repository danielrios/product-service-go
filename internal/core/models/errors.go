package models

import "errors"

// Erros comuns de domínio para Product.
var (
	ErrProductNotFound      = errors.New("product not found")
	ErrInvalidProductID     = errors.New("invalid product ID")
	ErrProductAlreadyExists = errors.New("product with this ID already exists")
)
