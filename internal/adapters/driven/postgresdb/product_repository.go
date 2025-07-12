package postgresdb

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/danielrios/product-service-go/internal/core/models"
	"github.com/danielrios/product-service-go/internal/core/ports"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // Importa o driver para registrar no database/sql
)

// PostgresProductRepository é a implementação do repositório para PostgreSQL.
type PostgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository cria uma nova instância do repositório, conectando ao banco.
func NewPostgresProductRepository(dataSourceName string) (*PostgresProductRepository, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Verifica se a conexão com o banco de dados está realmente funcionando.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso.")
	return &PostgresProductRepository{db: db}, nil
}

// Garante em tempo de compilação que PostgresProductRepository implementa a interface.
var _ ports.ProductRepository = (*PostgresProductRepository)(nil)

// Add adiciona um novo produto ao banco de dados.
func (r *PostgresProductRepository) Add(product *models.Product) error {
	query := "INSERT INTO products (id, name, price, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(context.Background(), query, product.ID, product.Name, product.Price, product.CreatedAt)

	if err != nil {
		// Verifica se o erro é de violação de chave única (produto já existe).
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 é o código para unique_violation
			return models.ErrProductAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID busca um produto pelo seu ID no banco de dados.
func (r *PostgresProductRepository) GetByID(id string) (*models.Product, error) {
	query := "SELECT id, name, price, created_at FROM products WHERE id = $1"
	row := r.db.QueryRowContext(context.Background(), query, id)

	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

// GetAll busca todos os produtos no banco de dados.
func (r *PostgresProductRepository) GetAll() ([]*models.Product, error) {
	query := "SELECT id, name, price, created_at FROM products"
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	// Verifica se houve algum erro durante a iteração das linhas.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Update atualiza um produto existente no banco de dados.
func (r *PostgresProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2 WHERE id = $3"
	result, err := r.db.ExecContext(context.Background(), query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

// Delete remove um produto do banco de dados pelo seu ID.
func (r *PostgresProductRepository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

func (r *PostgresProductRepository) Close() error {
	return r.db.Close()
}
