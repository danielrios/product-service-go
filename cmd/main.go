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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	httpDriver "github.com/danielrios/product-service-go/internal/adapters/driver/http"

	"github.com/danielrios/product-service-go/internal/adapters/driven/postgresdb"
	"github.com/danielrios/product-service-go/internal/application"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente do sistema.")
	}

	log.Println("Iniciando o microsserviço de produtos com Arquitetura Hexagonal...")

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
	productHandler := httpDriver.NewProductHandler(productService)

	// --- 4. Configura as Rotas HTTP com chi ---
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAllProductsHandler)
		r.Post("/", productHandler.CreateProductHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetProductByIDHandler)
			r.Put("/", productHandler.UpdateProductHandler)
			r.Delete("/", productHandler.DeleteProductHandler)
		})
	})

	// --- 5. Inicia o Servidor HTTP ---
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Servidor escutando em %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao escutar na porta %s: %v\n", server.Addr, err)
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
