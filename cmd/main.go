package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	http2 "github.com/danielrios/product-service-go/internal/adapters/driver/http"

	"github.com/danielrios/product-service-go/internal/adapters/driven/postgresdb"
	"github.com/danielrios/product-service-go/internal/application"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env.
	// É seguro ignorar o erro aqui, pois em ambientes de produção (como Docker ou Kubernetes),
	// as variáveis podem ser injetadas diretamente, e o arquivo .env pode não existir.
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente do sistema.")
	}

	log.Println("Iniciando o microsserviço de produtos com Arquitetura Hexagonal...")

	// --- Obter a string de conexão do banco de dados (idealmente de variáveis de ambiente) ---
	// Exemplo: "postgres://user:password@localhost:5432/database_name?sslmode=disable"
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	if dbConnectionString == "" {
		log.Fatal("A variável de ambiente DB_CONNECTION_STRING não está definida.")
	}

	// --- 1. Inicializa o Driven Adapter (Repositório de Produtos) ---
	productRepo, err := postgresdb.NewPostgresProductRepository(dbConnectionString)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	// --- 2. Inicializa o Application Service (Core) ---
	productService := application.NewProductService(productRepo)

	// --- 3. Inicializa o Driving Adapter (Handler HTTP) ---
	productHandler := http2.NewProductHandler(productService)

	// --- 4. Configura as Rotas HTTP ---
	mux := http.NewServeMux()

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetAllProductsHandler(w, r)
		case http.MethodPost:
			productHandler.CreateProductHandler(w, r)
		default:
			http.Error(w, "Método Não Permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/products" || r.URL.Path == "/products/" {
			mux.ServeHTTP(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductByIDHandler(w, r)
		case http.MethodPut:
			productHandler.UpdateProductHandler(w, r)
		case http.MethodDelete:
			productHandler.DeleteProductHandler(w, r)
		default:
			http.Error(w, "Método Não Permitido", http.StatusMethodNotAllowed)
		}
	})

	// --- 5. Inicia o Servidor HTTP ---
	port := ":8080"
	log.Printf("Microsserviço de produtos iniciado na porta %s", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Cria um canal para receber sinais do sistema operacional
	quit := make(chan os.Signal, 1)
	// Registra para receber sinais de interrupção (Ctrl+C) e término (SIGTERM)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine separada para não bloquear o main.
	go func() {
		log.Printf("Servidor escutando em %s...", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao escutar na porta %s: %v\n", port, err)
		}
	}()

	<-quit
	log.Println("Sinal de encerramento recebido. Desligando o servidor...")

	// --- 6. Graceful Shutdown ---
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Servidor forçado a desligar: %v", err)
	}

	log.Println("Servidor desligado graciosamente.")
}
