package helpers

import (
	"boilerplate-go/data/protocols"
	"github.com/go-redis/redis"
)

type RedisDatabaseIsConnectedHelper struct {
	client *redis.Client
}

func NewRedisDatabaseIsConnectedHelper(client *redis.Client) protocols.DatabaseIsConnected {
	return RedisDatabaseIsConnectedHelper{client}
}

func (helper RedisDatabaseIsConnectedHelper) IsConnected() (bool, error) {
	_, err := helper.client.Ping().Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
