package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderStatus representa os possíveis estados de um pedido.
type OrderStatus string

const (
	// StatusCreated indica que o pedido foi criado mas ainda não processado.
	StatusCreated OrderStatus = "CRIADO"
	// StatusProcessing indica que o pedido está sendo processado.
	StatusProcessing OrderStatus = "PROCESSANDO"
	// StatusProcessed indica que o pedido foi processado com sucesso.
	StatusProcessed OrderStatus = "PROCESSADO"
)

// Order é a entidade central do domínio.
// Representa um pedido no sistema.
type Order struct {
	id        string
	product   string
	quantity  int
	status    OrderStatus
	createdAt time.Time
}

// Erros de domínio
var (
	ErrInvalidProduct  = errors.New("product name cannot be empty")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)

// NewOrder cria uma nova ordem com validações de domínio.
func NewOrder(product string, quantity int) (*Order, error) {
	if product == "" {
		return nil, ErrInvalidProduct
	}

	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	return &Order{
		id:        uuid.New().String(),
		product:   product,
		quantity:  quantity,
		status:    StatusCreated,
		createdAt: time.Now(),
	}, nil
}

// Reconstruct reconstrói uma ordem a partir de dados persistidos.
// Usado pelo repository ao carregar do banco.
func Reconstruct(id, product string, quantity int, status OrderStatus, createdAt int64) *Order {
	return &Order{
		id:        id,
		product:   product,
		quantity:  quantity,
		status:    status,
		createdAt: time.Unix(createdAt, 0),
	}
}

// ID retorna o identificador único do pedido.
func (o *Order) ID() string {
	return o.id
}

// Product retorna o nome do produto.
func (o *Order) Product() string {
	return o.product
}

// Quantity retorna a quantidade do pedido.
func (o *Order) Quantity() int {
	return o.quantity
}

// Status retorna o status atual do pedido.
func (o *Order) Status() OrderStatus {
	return o.status
}

// CreatedAt retorna a data de criação do pedido.
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// MarkAsProcessing marca o pedido como em processamento.
func (o *Order) MarkAsProcessing() {
	o.status = StatusProcessing
}

// MarkAsProcessed marca o pedido como processado.
func (o *Order) MarkAsProcessed() {
	o.status = StatusProcessed
}
