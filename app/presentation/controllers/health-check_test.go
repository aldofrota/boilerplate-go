package controllers_test

import (
	"net/http"
	"testing"

	"boilerplate-go/app/presentation/controllers"
	"github.com/stretchr/testify/assert"
)

type ValidateIfHealthyUseCaseSpy struct {
	count       int
	result      bool
	err         error
	returnFalse bool
	returnError bool
}

func (spy *ValidateIfHealthyUseCaseSpy) Validate() (bool, error) {
	spy.count++
	spy.err = nil
	spy.result = true
	if spy.returnError {
		spy.err = assert.AnError
		return false, spy.err
	}
	if spy.returnFalse {
		return false, spy.err
	}
	return spy.result, spy.err
}

func TestHealthCheckController_Handle(t *testing.T) {
	t.Run("Should return 500 if use case returns an error", func(t *testing.T) {
		validateIfHealthyUseCaseSpy := &ValidateIfHealthyUseCaseSpy{returnError: true}
		sut := controllers.NewHealthCheckController(validateIfHealthyUseCaseSpy)
		result := sut.Handle(nil)
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
		assert.Equal(t, assert.AnError.Error(), result.Body)
		assert.Equal(t, 1, validateIfHealthyUseCaseSpy.count)
	})
	t.Run("Should return 500 if use case returns false", func(t *testing.T) {
		validateIfHealthyUseCaseSpy := &ValidateIfHealthyUseCaseSpy{returnFalse: true}
		sut := controllers.NewHealthCheckController(validateIfHealthyUseCaseSpy)
		result := sut.Handle(nil)
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
		assert.Equal(t, map[string]string{"message": "Service is not healthy"}, result.Body)
		assert.Equal(t, 1, validateIfHealthyUseCaseSpy.count)
	})
	t.Run("Should return 200 if use case returns true", func(t *testing.T) {
		validateIfHealthyUseCaseSpy := &ValidateIfHealthyUseCaseSpy{}
		sut := controllers.NewHealthCheckController(validateIfHealthyUseCaseSpy)
		result := sut.Handle(nil)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, map[string]string{"message": "Service is healthy"}, result.Body)
		assert.Equal(t, 1, validateIfHealthyUseCaseSpy.count)
	})
}
