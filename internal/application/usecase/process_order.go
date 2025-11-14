package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
	"github.com/diiineeei/neoway-backend/internal/domain/port"
)

// ProcessOrderUseCase é o caso de uso para processamento de pedidos.
type ProcessOrderUseCase struct {
	orderRepo port.OrderRepository
}

// NewProcessOrderUseCase cria uma nova instância do caso de uso.
func NewProcessOrderUseCase(repo port.OrderRepository) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{
		orderRepo: repo,
	}
}

// Execute executa o processamento de um pedido.
func (uc *ProcessOrderUseCase) Execute(ctx context.Context, orderID string) error {
	// 1. Simula tempo de processamento (requisito do teste)
	time.Sleep(2 * time.Second)

	// 2. Atualiza o status do pedido
	if err := uc.orderRepo.UpdateStatus(ctx, orderID, entity.StatusProcessed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
