package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
	"github.com/diiineeei/neoway-backend/internal/domain/port"
)

// MockOrderRepository é um mock do repositório para testes
type MockOrderRepository struct {
	SaveFunc         func(ctx context.Context, order *entity.Order) error
	FindByIDFunc     func(ctx context.Context, orderID string) (*entity.Order, error)
	UpdateStatusFunc func(ctx context.Context, orderID string, status entity.OrderStatus) error
}

func (m *MockOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) FindByID(ctx context.Context, orderID string) (*entity.Order, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, orderID)
	}
	return nil, nil
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, orderID, status)
	}
	return nil
}

// MockMessagePublisher é um mock do publisher para testes
type MockMessagePublisher struct {
	PublishFunc func(ctx context.Context, msg port.OrderMessage) error
}

func (m *MockMessagePublisher) PublishOrderMessage(ctx context.Context, msg port.OrderMessage) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, msg)
	}
	return nil
}

// TestCreateOrderUseCase_Success testa a criação bem-sucedida de um pedido
func TestCreateOrderUseCase_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{
		SaveFunc: func(ctx context.Context, order *entity.Order) error {
			return nil // Simula sucesso ao salvar
		},
	}

	mockPublisher := &MockMessagePublisher{
		PublishFunc: func(ctx context.Context, msg port.OrderMessage) error {
			return nil // Simula sucesso ao publicar
		},
	}

	useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

	input := CreateOrderInput{
		Product:  "Notebook Dell",
		Quantity: 2,
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	if err != nil {
		t.Errorf("Não esperava erro, mas recebeu: %v", err)
	}

	if output == nil {
		t.Fatal("Esperava um output válido, mas recebeu nil")
	}

	if output.OrderID == "" {
		t.Error("Esperava um OrderID não vazio")
	}

	if output.Status != entity.StatusCreated {
		t.Errorf("Esperava status %s, mas recebeu %s", entity.StatusCreated, output.Status)
	}
}

// TestCreateOrderUseCase_InvalidProduct testa erro com produto inválido
func TestCreateOrderUseCase_InvalidProduct(t *testing.T) {
	// Arrange
	mockRepo := &MockOrderRepository{}
	mockPublisher := &MockMessagePublisher{}

	useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

	input := CreateOrderInput{
		Product:  "", // Produto vazio
		Quantity: 1,
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Esperava um erro, mas não recebeu nenhum")
	}

	if output != nil {
		t.Error("Esperava output nil quando há erro")
	}
}

// TestCreateOrderUseCase_RepositoryError testa erro ao salvar no repositório
func TestCreateOrderUseCase_RepositoryError(t *testing.T) {
	// Arrange
	expectedError := errors.New("database connection error")

	mockRepo := &MockOrderRepository{
		SaveFunc: func(ctx context.Context, order *entity.Order) error {
			return expectedError // Simula erro no banco
		},
	}

	mockPublisher := &MockMessagePublisher{}

	useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

	input := CreateOrderInput{
		Product:  "Mouse Gamer",
		Quantity: 1,
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Esperava um erro, mas não recebeu nenhum")
	}

	if output != nil {
		t.Error("Esperava output nil quando há erro")
	}
}

// TestCreateOrderUseCase_PublisherError testa erro ao publicar mensagem
func TestCreateOrderUseCase_PublisherError(t *testing.T) {
	// Arrange
	expectedError := errors.New("rabbitmq connection error")

	mockRepo := &MockOrderRepository{
		SaveFunc: func(ctx context.Context, order *entity.Order) error {
			return nil // Sucesso ao salvar
		},
	}

	mockPublisher := &MockMessagePublisher{
		PublishFunc: func(ctx context.Context, msg port.OrderMessage) error {
			return expectedError // Simula erro ao publicar
		},
	}

	useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

	input := CreateOrderInput{
		Product:  "Teclado Mecânico",
		Quantity: 1,
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Esperava um erro, mas não recebeu nenhum")
	}

	if output != nil {
		t.Error("Esperava output nil quando há erro")
	}
}

// TestCreateOrderUseCase_MessageContent testa o conteúdo da mensagem publicada
func TestCreateOrderUseCase_MessageContent(t *testing.T) {
	// Arrange
	var capturedMessage port.OrderMessage

	mockRepo := &MockOrderRepository{
		SaveFunc: func(ctx context.Context, order *entity.Order) error {
			return nil
		},
	}

	mockPublisher := &MockMessagePublisher{
		PublishFunc: func(ctx context.Context, msg port.OrderMessage) error {
			capturedMessage = msg // Captura a mensagem publicada
			return nil
		},
	}

	useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

	input := CreateOrderInput{
		Product:  "Webcam",
		Quantity: 5,
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	if err != nil {
		t.Fatalf("Não esperava erro, mas recebeu: %v", err)
	}

	if capturedMessage.OrderID != output.OrderID {
		t.Errorf("OrderID na mensagem: esperava '%s', recebeu '%s'", output.OrderID, capturedMessage.OrderID)
	}

	if capturedMessage.Status != entity.StatusProcessing {
		t.Errorf("Status na mensagem: esperava %s, recebeu %s", entity.StatusProcessing, capturedMessage.Status)
	}
}

// TestCreateOrderUseCase_TableDriven testa múltiplos cenários
func TestCreateOrderUseCase_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateOrderInput
		repoError     error
		publishError  error
		shouldSucceed bool
	}{
		{
			name:          "Sucesso - Notebook",
			input:         CreateOrderInput{Product: "Notebook", Quantity: 1},
			repoError:     nil,
			publishError:  nil,
			shouldSucceed: true,
		},
		{
			name:          "Sucesso - múltiplos itens",
			input:         CreateOrderInput{Product: "Mouse", Quantity: 10},
			repoError:     nil,
			publishError:  nil,
			shouldSucceed: true,
		},
		{
			name:          "Erro - produto vazio",
			input:         CreateOrderInput{Product: "", Quantity: 1},
			repoError:     nil,
			publishError:  nil,
			shouldSucceed: false,
		},
		{
			name:          "Erro - quantidade zero",
			input:         CreateOrderInput{Product: "Teclado", Quantity: 0},
			repoError:     nil,
			publishError:  nil,
			shouldSucceed: false,
		},
		{
			name:          "Erro - falha no repositório",
			input:         CreateOrderInput{Product: "Monitor", Quantity: 1},
			repoError:     errors.New("db error"),
			publishError:  nil,
			shouldSucceed: false,
		},
		{
			name:          "Erro - falha ao publicar",
			input:         CreateOrderInput{Product: "Webcam", Quantity: 1},
			repoError:     nil,
			publishError:  errors.New("queue error"),
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockOrderRepository{
				SaveFunc: func(ctx context.Context, order *entity.Order) error {
					return tt.repoError
				},
			}

			mockPublisher := &MockMessagePublisher{
				PublishFunc: func(ctx context.Context, msg port.OrderMessage) error {
					return tt.publishError
				},
			}

			useCase := NewCreateOrderUseCase(mockRepo, mockPublisher)

			// Act
			output, err := useCase.Execute(context.Background(), tt.input)

			// Assert
			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Não esperava erro, mas recebeu: %v", err)
				}
				if output == nil {
					t.Error("Esperava output válido")
				}
			} else {
				if err == nil {
					t.Error("Esperava um erro")
				}
				if output != nil {
					t.Error("Esperava output nil quando há erro")
				}
			}
		})
	}
}
