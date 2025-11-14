package handler

import (
	"net/http"

	"github.com/diiineeei/neoway-backend/internal/adapter/input/http/dto"
	"github.com/gin-gonic/gin"
)

// HealthHandler gerencia requisições de health check.
type HealthHandler struct{}

// NewHealthHandler cria uma nova instância do health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check godoc
// @Summary      Health Check
// @Description  Verifica se a API está funcionando
// @Tags         health
// @Produce      json
// @Success      200 {object} dto.HealthResponse "API está saudável"
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  "healthy",
		Service: "api",
	})
}
