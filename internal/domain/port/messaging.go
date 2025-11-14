package port

import (
	"context"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
)

// OrderMessage representa uma mensagem de pedido.
type OrderMessage struct {
	OrderID string              `json:"order_id"`
	Status  entity.OrderStatus  `json:"status"`
}

// MessagePublisher define a interface (porta) para publicação de mensagens.
// Esta é uma porta de saída (output port) na arquitetura hexagonal.
type MessagePublisher interface {
	// PublishOrderMessage publica uma mensagem de pedido.
	PublishOrderMessage(ctx context.Context, msg OrderMessage) error
}

// MessageConsumer define a interface (porta) para consumo de mensagens.
// Esta é uma porta de saída (output port) na arquitetura hexagonal.
type MessageConsumer interface {
	// ConsumeOrderMessages consome mensagens de pedidos.
	ConsumeOrderMessages(ctx context.Context, handler func(OrderMessage) error) error
}
