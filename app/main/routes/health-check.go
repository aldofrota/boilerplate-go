package routes

import (
	"boilerplate-go/app/main/adapters"
	"boilerplate-go/app/main/factories"
	"github.com/gin-gonic/gin"
)

func addHealthCheckRoutes(rg *gin.RouterGroup) {
	healthCheck := rg.Group("/health")
	healthCheck.GET("", adapters.AdaptController(factories.NewHealthCheckControllerFactory()))
}
