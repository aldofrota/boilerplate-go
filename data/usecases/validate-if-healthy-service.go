package usecases

import (
	"errors"
	"sync"

	"boilerplate-go/data/protocols"
	"boilerplate-go/domain/usecases"
)

// ValidateIfHealthy is a use case that validates if the service is healthy
type ValidateIfHealthyService struct {
	redisDatabaseIsConnected protocols.DatabaseIsConnected
	mongoDatabaseIsConnected protocols.DatabaseIsConnected
}

func NewValidateIfHealthyService(
	redisDatabaseIsConnected protocols.DatabaseIsConnected,
	mongoDatabaseIsConnected protocols.DatabaseIsConnected,
) usecases.ValidateIfHealthy {
	return ValidateIfHealthyService{
		redisDatabaseIsConnected,
		mongoDatabaseIsConnected,
	}
}

func (service ValidateIfHealthyService) Validate() (bool, error) {
	var wg sync.WaitGroup
	count := 0
	wg.Add(2)

	resultChannel := make(chan bool)
	errChannel := make(chan error, 1)

	checkConnection := func(isConnectedFunc func() (bool, error)) {
		defer wg.Done()
		connected, err := isConnectedFunc()
		if err != nil || !connected {
			errChannel <- err
			return
		}
		resultChannel <- true
	}

	go checkConnection(service.redisDatabaseIsConnected.IsConnected)
	go checkConnection(service.mongoDatabaseIsConnected.IsConnected)

	go func() {
		wg.Wait()
		close(resultChannel)
		close(errChannel)
	}()

	for result := range resultChannel {
		count++
		if !result {
			return false, <-errChannel
		}
	}

	if count != 2 {
		return false, errors.New("failed to check all services")
	}

	return true, nil
}
