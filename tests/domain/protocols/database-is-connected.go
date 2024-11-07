package protocols

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase simula a interface DatabaseIsConnected para testes
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) IsConnected() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

// Test cases for DatabaseIsConnected interface
func TestDatabaseIsConnected(t *testing.T) {
	mockDB := new(MockDatabase)

	// Cenário: Banco de dados conectado
	mockDB.On("IsConnected").Return(true, nil)

	connected, err := mockDB.IsConnected()
	assert.NoError(t, err)
	assert.True(t, connected)

	// Cenário: Banco de dados desconectado
	mockDB.On("IsConnected").Return(false, nil)

	connected, err = mockDB.IsConnected()
	assert.NoError(t, err)
	assert.False(t, connected)

	// Cenário: Erro ao conectar no banco
	mockDB.On("IsConnected").Return(false, assert.AnError)

	connected, err = mockDB.IsConnected()
	assert.Error(t, err)
	assert.False(t, connected)

	// Verificar se todos os métodos mockados foram chamados
	mockDB.AssertExpectations(t)
}
