package routes

import (
	"boilerplate-go/main/adapters"
	"boilerplate-go/main/factories"
	"github.com/gin-gonic/gin"
)

func addHealthCheckRoutes(rg *gin.RouterGroup) {
	healthCheck := rg.Group("/health")
	healthCheck.GET("", adapters.AdaptController(factories.NewHealthCheckControllerFactory()))
}
