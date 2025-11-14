package dto

import "github.com/diiineeei/neoway-backend/internal/domain/entity"

// CreateOrderRequest requisição para criar um pedido
type CreateOrderRequest struct {
	Product  string `json:"product" binding:"required" example:"Notebook Dell"`
	Quantity int    `json:"quantity" binding:"required,min=1" example:"2"`
}

// CreateOrderResponse resposta ao criar um pedido
type CreateOrderResponse struct {
	OrderID string             `json:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status  entity.OrderStatus `json:"status" example:"CRIADO"`
}

// ErrorResponse resposta de erro
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Details string `json:"details,omitempty" example:"product is required"`
}

// HealthResponse resposta do health check
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"api"`
}
