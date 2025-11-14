package usecase

import (
	"context"
	"fmt"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
	"github.com/diiineeei/neoway-backend/internal/domain/port"
)

// CreateOrderInput representa os dados de entrada para criar um pedido.
type CreateOrderInput struct {
	Product  string
	Quantity int
}

// CreateOrderOutput representa o resultado da criação de um pedido.
type CreateOrderOutput struct {
	OrderID string
	Status  entity.OrderStatus
}

// CreateOrderUseCase é o caso de uso para criação de pedidos.
type CreateOrderUseCase struct {
	orderRepo port.OrderRepository
	publisher port.MessagePublisher
}

// NewCreateOrderUseCase cria uma nova instância do caso de uso.
func NewCreateOrderUseCase(repo port.OrderRepository, pub port.MessagePublisher) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo: repo,
		publisher: pub,
	}
}

// Execute executa o caso de uso de criar pedido.
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// 1. Cria a entidade de domínio (com validações)
	order, err := entity.NewOrder(input.Product, input.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create order entity: %w", err)
	}

	// 2. Persiste o pedido
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// 3. Publica mensagem para processamento assíncrono
	message := port.OrderMessage{
		OrderID: order.ID(),
		Status:  entity.StatusProcessing,
	}

	if err := uc.publisher.PublishOrderMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// 4. Retorna o resultado
	return &CreateOrderOutput{
		OrderID: order.ID(),
		Status:  order.Status(),
	}, nil
}
