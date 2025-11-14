package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/diiineeei/neoway-backend/docs"
	"github.com/diiineeei/neoway-backend/internal/adapter/input/http/handler"
	"github.com/diiineeei/neoway-backend/internal/adapter/input/http/router"
	"github.com/diiineeei/neoway-backend/internal/adapter/output/kafka"
	"github.com/diiineeei/neoway-backend/internal/adapter/output/mongodb"
	"github.com/diiineeei/neoway-backend/internal/application/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           Neoway Backend API
// @version         1.0
// @description     Sistema de processamento assíncrono de pedidos
// @host            localhost:8080
// @BasePath        /
func main() {
	mongoURI := getEnv("MONGODB_URI", "mongodb://mongodb:27017")
	mongoDatabase := getEnv("MONGODB_DATABASE", "neoway")
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "orders")
	port := getEnv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Adapters (Infraestrutura)
	orderCollection := mongoClient.Database(mongoDatabase).Collection("orders")
	orderRepository := mongodb.NewOrderRepositoryAdapterFromCollection(orderCollection)

	// Inicializa Kafka Publisher
	brokers := strings.Split(kafkaBrokers, ",")
	kafkaPublisher := kafka.NewMessagePublisherAdapter(kafka.Config{
		Brokers: brokers,
		Topic:   kafkaTopic,
	})
	defer kafkaPublisher.Close()
	log.Println("Successfully connected to Kafka")

	// Use Cases (Camada de Aplicação)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, kafkaPublisher)

	// Handlers (Adaptadores de Entrada)
	orderHandler := handler.NewOrderHandler(createOrderUseCase)
	healthHandler := handler.NewHealthHandler()

	// Configura o Gin
	gin.SetMode(gin.ReleaseMode)

	// Setup Router com arquitetura hexagonal
	r := router.SetupRouter(router.Config{
		OrderHandler:  orderHandler,
		HealthHandler: healthHandler,
	})

	// Configura servidor HTTP
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Inicia servidor em goroutine
	go func() {
		log.Printf("Starting API server on port %s", port)
		log.Printf("Swagger UI available at http://localhost:%s/doc", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguarda sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
