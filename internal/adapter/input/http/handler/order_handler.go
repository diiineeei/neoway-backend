package handler

import (
	"net/http"

	"github.com/diiineeei/neoway-backend/internal/adapter/input/http/dto"
	"github.com/diiineeei/neoway-backend/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

// OrderHandler gerencia as requisições HTTP relacionadas a pedidos.
type OrderHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
}

// NewOrderHandler cria uma nova instância do handler.
func NewOrderHandler(createOrderUseCase *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
	}
}

// CreateOrder godoc
// @Summary      Criar novo pedido
// @Description  Cria um novo pedido e envia para processamento assíncrono
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderRequest true "Dados do pedido"
// @Success      201 {object} dto.CreateOrderResponse "Pedido criado com sucesso"
// @Failure      400 {object} dto.ErrorResponse "Dados inválidos"
// @Failure      500 {object} dto.ErrorResponse "Erro interno do servidor"
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	input := usecase.CreateOrderInput{
		Product:  req.Product,
		Quantity: req.Quantity,
	}

	output, err := h.createOrderUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create order",
		})
		return
	}

	response := dto.CreateOrderResponse{
		OrderID: output.OrderID,
		Status:  output.Status,
	}

	c.JSON(http.StatusCreated, response)
}
