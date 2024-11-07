package adapters

import (
	"boilerplate-go/app/presentation/protocols"
	"github.com/gin-gonic/gin"
)

// AdaptController adapts a controller to a gin handler
// This is a function that receives a controller and returns a gin.HandlerFunc
func AdaptController(controller protocols.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpResponse := controller.Handle(c)
		if httpResponse.StatusCode >= 200 && httpResponse.StatusCode <= 299 {
			c.JSON(httpResponse.StatusCode, httpResponse.Body)
		} else {
			c.JSON(httpResponse.StatusCode, gin.H{"error": httpResponse.Body})
		}
	}
}
