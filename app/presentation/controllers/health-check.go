package controllers

import (
	"net/http"

	"boilerplate-go/app/domain/usecases"
	"boilerplate-go/app/presentation/protocols"
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	validateIfHealthyUseCase usecases.ValidateIfHealthy
}

func NewHealthCheckController(validateIfHealthyUseCase usecases.ValidateIfHealthy) protocols.Controller {
	return &HealthCheckController{validateIfHealthyUseCase}
}

func (controller *HealthCheckController) Handle(ctx *gin.Context) protocols.HttpResponse {
	result, err := controller.validateIfHealthyUseCase.Validate()
	if err != nil {
		return protocols.HttpResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}
	}
	if !result {
		return protocols.HttpResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       map[string]string{"message": "Service is not healthy"},
		}
	}
	return protocols.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       map[string]string{"message": "Service is healthy"},
	}
}
