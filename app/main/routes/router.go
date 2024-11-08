package routes

import (
	"context"
	"net/http"
	"os"

	"boilerplate-go/app/main/middlewares"
	"github.com/gin-gonic/gin"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

var router = gin.New()
var server = &http.Server{
	Addr:    ":3000",
	Handler: router,
}

func Run() error {
	getRoutes()
	port := os.Getenv("PORT")
	if port != "" {
		server.Addr = ":" + port
	}
	return server.ListenAndServe()
}

func ShutDown(ctx context.Context) error {
	return server.Shutdown(ctx)
}

func getRoutes() {
	router.Use(gin.Recovery())
	router.Use(httptrace.Middleware(os.Getenv("DD_SERVICE")))
	router.Use(middlewares.NewCorsMiddleware())

	apiPrefix := router.Group("")
	addDocsRoutes(apiPrefix)
	addHealthCheckRoutes(apiPrefix)
}
