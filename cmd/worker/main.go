package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diiineeei/neoway-backend/internal/adapter/output/kafka"
	"github.com/diiineeei/neoway-backend/internal/adapter/output/mongodb"
	"github.com/diiineeei/neoway-backend/internal/application/usecase"
	"github.com/diiineeei/neoway-backend/internal/domain/port"
	"github.com/diiineeei/neoway-backend/pkg/env"
)

func main() {
	// Configurações a partir de variáveis de ambiente
	mongoURI := env.GetString("MONGODB_URI", "mongodb://mongodb:27017")
	mongoDatabase := env.GetString("MONGODB_DATABASE", "neoway")
	kafkaTopic := env.GetString("KAFKA_TOPIC", "orders")
	kafkaGroupID := env.GetString("KAFKA_GROUP_ID", "neoway-workers")

	// Inicializa MongoDB
	mongoConn, err := mongodb.NewConnection(mongodb.Config{
		URI:      mongoURI,
		Database: mongoDatabase,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoConn.Close(context.Background())

	log.Println("✓ Successfully connected to MongoDB")

	// Inicializa Kafka Consumer
	brokers := env.GetStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}, ",")
	kafkaConsumer := kafka.NewMessageConsumerAdapter(kafka.Config{
		Brokers: brokers,
		Topic:   kafkaTopic,
		GroupID: kafkaGroupID,
	})
	defer kafkaConsumer.Close()

	log.Println("✓ Successfully connected to Kafka")

	// Cria adapters de saída (output ports)
	orderRepo := mongodb.NewOrderRepositoryAdapter(mongoConn.Database())

	// Cria use case (application layer)
	processOrderUseCase := usecase.NewProcessOrderUseCase(orderRepo)

	// Cria contexto com cancelamento
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Canal para erros do worker
	errChan := make(chan error, 1)

	// Inicia worker em goroutine
	go func() {
		log.Println("🔄 Starting worker to consume messages...")

		handler := createMessageHandler(processOrderUseCase)

		if err := kafkaConsumer.ConsumeOrderMessages(ctx, handler); err != nil {
			errChan <- fmt.Errorf("worker error: %w", err)
		}
	}()

	// Aguarda sinal de interrupção ou erro
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("🛑 Received shutdown signal...")
	case err := <-errChan:
		log.Printf("❌ Worker error: %v", err)
	}

	// Cancela contexto para parar consumo
	cancel()

	log.Println("✓ Worker stopped")
}

// createMessageHandler cria o handler para processar mensagens.
func createMessageHandler(processOrderUseCase *usecase.ProcessOrderUseCase) func(port.OrderMessage) error {
	return func(msg port.OrderMessage) error {
		log.Printf("📦 Processing order: %s", msg.OrderID)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := processOrderUseCase.Execute(ctx, msg.OrderID); err != nil {
			log.Printf("❌ Failed to process order %s: %v", msg.OrderID, err)
			return err
		}

		log.Printf("✓ Order %s processed successfully", msg.OrderID)
		return nil
	}
}
