package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/diiineeei/neoway-backend/internal/domain/port"
	"github.com/segmentio/kafka-go"
)

// MessageConsumerAdapter é o adaptador de saída para consumo de mensagens do Kafka.
// Implementa port.MessageConsumer.
type MessageConsumerAdapter struct {
	reader *kafka.Reader
}

// NewMessageConsumerAdapter cria um novo adapter de consumer.
func NewMessageConsumerAdapter(config Config) *MessageConsumerAdapter {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		Topic:       config.Topic,
		GroupID:     config.GroupID,
		MinBytes:    10e1, // 100B
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})

	return &MessageConsumerAdapter{
		reader: reader,
	}
}

// ConsumeOrderMessages consome mensagens de pedidos do Kafka.
func (c *MessageConsumerAdapter) ConsumeOrderMessages(ctx context.Context, handler func(port.OrderMessage) error) error {
	log.Println("Starting Kafka consumer...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping consumer")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Printf("Error fetching message: %v", err)
				continue
			}

			var orderMsg port.OrderMessage
			if err := json.Unmarshal(msg.Value, &orderMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				// Commit mesmo com erro para não reprocessar mensagem inválida
				if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
					log.Printf("Error committing invalid message: %v", commitErr)
				}
				continue
			}

			log.Printf("Received message: OrderID=%s, Status=%s", orderMsg.OrderID, orderMsg.Status)

			// Processa a mensagem
			if err := handler(orderMsg); err != nil {
				log.Printf("Error processing message: %v", err)
				// Não faz commit, permitindo reprocessamento
				continue
			}

			// Commit da mensagem processada com sucesso
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

// Close fecha a conexão com o Kafka.
func (c *MessageConsumerAdapter) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close kafka reader: %w", err)
	}
	return nil
}
