package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diiineeei/neoway-backend/internal/domain/port"
	"github.com/segmentio/kafka-go"
)

// MessagePublisherAdapter é o adaptador de saída para publicação no Kafka.
// Implementa port.MessagePublisher.
type MessagePublisherAdapter struct {
	writer *kafka.Writer
}

// NewMessagePublisherAdapter cria um novo adapter de publisher.
func NewMessagePublisherAdapter(config Config) *MessagePublisherAdapter {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireOne,
		Async:                  false,
	}

	return &MessagePublisherAdapter{
		writer: writer,
	}
}

// PublishOrderMessage publica uma mensagem de pedido no Kafka.
func (p *MessagePublisherAdapter) PublishOrderMessage(ctx context.Context, msg port.OrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(msg.OrderID),
		Value: body,
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to publish message to kafka: %w", err)
	}

	return nil
}

// Close fecha a conexão com o Kafka.
func (p *MessagePublisherAdapter) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("failed to close kafka writer: %w", err)
	}
	return nil
}
