package helpers

import (
	"context"

	"github.com/go-redis/redis/v8"
	"boilerplate-go/app/domain/protocols"
)

type RedisDatabaseIsConnectedHelper struct {
	client redis.UniversalClient
}

func NewRedisDatabaseIsConnectedHelper(client redis.UniversalClient) protocols.DatabaseIsConnected {
	return RedisDatabaseIsConnectedHelper{client}
}

func (helper RedisDatabaseIsConnectedHelper) IsConnected() (bool, error) {
	ctx := context.Background()
	_, err := helper.client.Ping(ctx).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
