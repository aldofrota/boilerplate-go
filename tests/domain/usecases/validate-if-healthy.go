package usecases

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockValidateIfHealthy simula a interface ValidateIfHealthy para testes
type MockValidateIfHealthy struct {
	mock.Mock
}

func (m *MockValidateIfHealthy) Validate() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

// TestValidateIfHealthy testa a interface ValidateIfHealthy em diferentes cenários
func TestValidateIfHealthy(t *testing.T) {
	mockValidator := new(MockValidateIfHealthy)

	// Cenário: Serviço está saudável
	mockValidator.On("Validate").Return(true, nil)

	isHealthy, err := mockValidator.Validate()
	assert.NoError(t, err)
	assert.True(t, isHealthy)

	// Cenário: Serviço está com problemas (não saudável)
	mockValidator.On("Validate").Return(false, nil)

	isHealthy, err = mockValidator.Validate()
	assert.NoError(t, err)
	assert.False(t, isHealthy)

	// Cenário: Erro ao validar o serviço
	mockValidator.On("Validate").Return(false, assert.AnError)

	isHealthy, err = mockValidator.Validate()
	assert.Error(t, err)
	assert.False(t, isHealthy)

	// Verificar se todos os métodos mockados foram chamados conforme esperado
	mockValidator.AssertExpectations(t)
}
