package routes

import (
	_ "boilerplate-go/app/main/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func addDocsRoutes(rg *gin.RouterGroup) {
	rg.GET("/docs", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
