package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
)

// TestProcessOrderUseCase_Success testa o processamento bem-sucedido
func TestProcessOrderUseCase_Success(t *testing.T) {
	var updatedStatus entity.OrderStatus
	var updatedOrderID string

	mockRepo := &MockOrderRepository{
		UpdateStatusFunc: func(ctx context.Context, orderID string, status entity.OrderStatus) error {
			updatedOrderID = orderID
			updatedStatus = status
			return nil
		},
	}

	useCase := NewProcessOrderUseCase(mockRepo)
	orderID := "test-order-123"

	start := time.Now()
	err := useCase.Execute(context.Background(), orderID)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Não esperava erro, mas recebeu: %v", err)
	}

	if updatedOrderID != orderID {
		t.Errorf("OrderID: esperava '%s', recebeu '%s'", orderID, updatedOrderID)
	}

	if updatedStatus != entity.StatusProcessed {
		t.Errorf("Status: esperava %s, recebeu %s", entity.StatusProcessed, updatedStatus)
	}

	// Verifica se o delay de 2 segundos foi aplicado
	if duration < 2*time.Second {
		t.Errorf("Esperava delay de pelo menos 2s, mas levou %v", duration)
	}
}

// TestProcessOrderUseCase_RepositoryError testa erro no repositório
func TestProcessOrderUseCase_RepositoryError(t *testing.T) {
	expectedError := errors.New("database error")

	mockRepo := &MockOrderRepository{
		UpdateStatusFunc: func(ctx context.Context, orderID string, status entity.OrderStatus) error {
			return expectedError
		},
	}

	useCase := NewProcessOrderUseCase(mockRepo)

	err := useCase.Execute(context.Background(), "order-123")

	if err == nil {
		t.Error("Esperava um erro, mas não recebeu nenhum")
	}
}

// TestProcessOrderUseCase_ContextCancellation testa cancelamento via context
func TestProcessOrderUseCase_ContextCancellation(t *testing.T) {
	// Este teste verifica que o context é passado corretamente
	var receivedContext context.Context

	mockRepo := &MockOrderRepository{
		UpdateStatusFunc: func(ctx context.Context, orderID string, status entity.OrderStatus) error {
			receivedContext = ctx
			return nil
		},
	}

	useCase := NewProcessOrderUseCase(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = useCase.Execute(ctx, "order-123")

	if receivedContext != ctx {
		t.Error("Context não foi passado corretamente para o repositório")
	}
}

// TestProcessOrderUseCase_MultipleOrders testa processamento de múltiplos pedidos
func TestProcessOrderUseCase_MultipleOrders(t *testing.T) {
	processedOrders := make(map[string]bool)

	mockRepo := &MockOrderRepository{
		UpdateStatusFunc: func(ctx context.Context, orderID string, status entity.OrderStatus) error {
			processedOrders[orderID] = true
			return nil
		},
	}

	useCase := NewProcessOrderUseCase(mockRepo)

	orderIDs := []string{"order-1", "order-2", "order-3"}

	for _, orderID := range orderIDs {
		err := useCase.Execute(context.Background(), orderID)
		if err != nil {
			t.Errorf("Erro ao processar %s: %v", orderID, err)
		}
	}

	for _, orderID := range orderIDs {
		if !processedOrders[orderID] {
			t.Errorf("Pedido %s não foi processado", orderID)
		}
	}
}
