package entity

import (
	"errors"
	"testing"
	"time"
)

// TestNewOrder_Success testa a criação bem-sucedida de um pedido
func TestNewOrder_Success(t *testing.T) {
	// Arrange (Preparar)
	product := "Notebook Dell"
	quantity := 2

	// Act (Executar)
	order, err := NewOrder(product, quantity)

	// Assert (Verificar)
	if err != nil {
		t.Errorf("Erro recebido: %v", err)
	}

	if order == nil {
		t.Fatal("Esperava uma order válida, mas recebeu nil")
	}

	if order.Product() != product {
		t.Errorf("Esperava product '%s', mas recebeu '%s'", product, order.Product())
	}

	if order.Quantity() != quantity {
		t.Errorf("Esperava quantity %d, mas recebeu %d", quantity, order.Quantity())
	}

	if order.Status() != StatusCreated {
		t.Errorf("Esperava status %s, mas recebeu %s", StatusCreated, order.Status())
	}

	if order.ID() == "" {
		t.Error("Esperava um ID não vazio")
	}

	if order.CreatedAt().IsZero() {
		t.Error("Esperava uma data de criação válida")
	}
}

// TestNewOrder_EmptyProduct testa erro quando produto está vazio
func TestNewOrder_EmptyProduct(t *testing.T) {
	// Arrange
	product := ""
	quantity := 1

	// Act
	order, err := NewOrder(product, quantity)

	// Assert
	if !errors.Is(err, ErrInvalidProduct) {
		t.Errorf("Esperava erro ErrInvalidProduct, mas recebeu: %v", err)
	}

	if order != nil {
		t.Error("Esperava order nil quando há erro")
	}
}

// TestNewOrder_ZeroQuantity testa erro quando quantidade é zero
func TestNewOrder_ZeroQuantity(t *testing.T) {
	// Arrange
	product := "Mouse Gamer"
	quantity := 0

	// Act
	order, err := NewOrder(product, quantity)

	// Assert
	if !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("Esperava erro ErrInvalidQuantity, mas recebeu: %v", err)
	}

	if order != nil {
		t.Error("Esperava order nil quando há erro")
	}
}

// TestNewOrder_NegativeQuantity testa erro quando quantidade é negativa
func TestNewOrder_NegativeQuantity(t *testing.T) {
	// Arrange
	product := "Teclado"
	quantity := -5

	// Act
	order, err := NewOrder(product, quantity)

	// Assert
	if !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("Esperava erro ErrInvalidQuantity, mas recebeu: %v", err)
	}

	if order != nil {
		t.Error("Esperava order nil quando há erro")
	}
}

// TestNewOrder_TableDriven testa múltiplos cenários usando table-driven tests
func TestNewOrder_TableDriven(t *testing.T) {
	// Define os casos de teste
	tests := []struct {
		name        string
		product     string
		quantity    int
		shouldError bool
		expectedErr error
	}{
		{
			name:        "Pedido válido - Notebook",
			product:     "Notebook",
			quantity:    1,
			shouldError: false,
			expectedErr: nil,
		},
		{
			name:        "Pedido válido - múltiplos itens",
			product:     "Mouse",
			quantity:    10,
			shouldError: false,
			expectedErr: nil,
		},
		{
			name:        "Erro - produto vazio",
			product:     "",
			quantity:    1,
			shouldError: true,
			expectedErr: ErrInvalidProduct,
		},
		{
			name:        "Erro - quantidade zero",
			product:     "Teclado",
			quantity:    0,
			shouldError: true,
			expectedErr: ErrInvalidQuantity,
		},
		{
			name:        "Erro - quantidade negativa",
			product:     "Monitor",
			quantity:    -1,
			shouldError: true,
			expectedErr: ErrInvalidQuantity,
		},
	}

	// Executa cada caso de teste
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			order, err := NewOrder(tt.product, tt.quantity)

			// Assert
			if tt.shouldError {
				if err == nil {
					t.Errorf("Esperava um erro, mas não recebeu nenhum")
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Esperava erro %v, mas recebeu %v", tt.expectedErr, err)
				}
				if order != nil {
					t.Error("Esperava order nil quando há erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas recebeu: %v", err)
				}
				if order == nil {
					t.Fatal("Esperava uma order válida, mas recebeu nil")
				}
				if order.Product() != tt.product {
					t.Errorf("Product: esperava '%s', recebeu '%s'", tt.product, order.Product())
				}
				if order.Quantity() != tt.quantity {
					t.Errorf("Quantity: esperava %d, recebeu %d", tt.quantity, order.Quantity())
				}
			}
		})
	}
}

// TestReconstruct testa a reconstrução de uma order do banco de dados
func TestReconstruct(t *testing.T) {
	// Arrange
	id := "123e4567-e89b-12d3-a456-426614174000"
	product := "Notebook"
	quantity := 2
	status := StatusProcessed
	createdAt := time.Now().Unix()

	// Act
	order := Reconstruct(id, product, quantity, status, createdAt)

	// Assert
	if order == nil {
		t.Fatal("Esperava uma order válida, mas recebeu nil")
	}

	if order.ID() != id {
		t.Errorf("ID: esperava '%s', recebeu '%s'", id, order.ID())
	}

	if order.Product() != product {
		t.Errorf("Product: esperava '%s', recebeu '%s'", product, order.Product())
	}

	if order.Quantity() != quantity {
		t.Errorf("Quantity: esperava %d, recebeu %d", quantity, order.Quantity())
	}

	if order.Status() != status {
		t.Errorf("Status: esperava %s, recebeu %s", status, order.Status())
	}

	expectedTime := time.Unix(createdAt, 0)
	if !order.CreatedAt().Equal(expectedTime) {
		t.Errorf("CreatedAt: esperava %v, recebeu %v", expectedTime, order.CreatedAt())
	}
}

// TestMarkAsProcessing testa a marcação de pedido como processando
func TestMarkAsProcessing(t *testing.T) {
	// Arrange
	order, _ := NewOrder("Produto", 1)
	initialStatus := order.Status()

	// Act
	order.MarkAsProcessing()

	// Assert
	if order.Status() != StatusProcessing {
		t.Errorf("Esperava status %s, mas recebeu %s", StatusProcessing, order.Status())
	}

	if initialStatus == StatusProcessing {
		t.Error("Status inicial não deveria ser StatusProcessing")
	}
}

// TestMarkAsProcessed testa a marcação de pedido como processado
func TestMarkAsProcessed(t *testing.T) {
	// Arrange
	order, _ := NewOrder("Produto", 1)
	order.MarkAsProcessing()

	// Act
	order.MarkAsProcessed()

	// Assert
	if order.Status() != StatusProcessed {
		t.Errorf("Esperava status %s, mas recebeu %s", StatusProcessed, order.Status())
	}
}

// TestOrderStatus_Constants testa as constantes de status
func TestOrderStatus_Constants(t *testing.T) {
	tests := []struct {
		status   OrderStatus
		expected string
	}{
		{StatusCreated, "CRIADO"},
		{StatusProcessing, "PROCESSANDO"},
		{StatusProcessed, "PROCESSADO"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Esperava '%s', mas recebeu '%s'", tt.expected, string(tt.status))
			}
		})
	}
}

// BenchmarkNewOrder benchmarks a criação de pedidos
func BenchmarkNewOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewOrder("Produto Teste", 1)
	}
}
