package port

import (
	"context"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
)

// OrderRepository define a interface (porta) para persistência de pedidos.
// Esta é uma porta de saída (output port) na arquitetura hexagonal.
type OrderRepository interface {
	// Save persiste um pedido.
	Save(ctx context.Context, order *entity.Order) error

	// FindByID busca um pedido pelo ID.
	FindByID(ctx context.Context, orderID string) (*entity.Order, error)

	// UpdateStatus atualiza o status de um pedido.
	UpdateStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
}
