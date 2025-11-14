package router

import (
	"net/http"

	"github.com/diiineeei/neoway-backend/internal/adapter/input/http/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config contém as dependências necessárias para configurar as rotas.
type Config struct {
	OrderHandler  *handler.OrderHandler
	HealthHandler *handler.HealthHandler
}

// SetupRouter configura todas as rotas da aplicação.
func SetupRouter(config Config) *gin.Engine {
	router := gin.Default()

	router.GET("/health", config.HealthHandler.Check)

	router.GET("/doc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/doc/index.html")
	})
	router.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", config.OrderHandler.CreateOrder)
		}
	}

	router.POST("/orders", config.OrderHandler.CreateOrder)

	return router
}
