package factories

import (
	"context"
	"errors"
	"os"

	"github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

var db_redis_con redis.UniversalClient

func NewDatabaseRedisOpenConnection() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		return errors.New("REDIS_HOST or REDIS_PORT not found in environment variables")
	}

	// Configuração do cliente Redis com tracing
	opts := &redis.Options{
		Addr: redisHost + ":" + redisPort,
		DB:   0,
	}

	// Use o cliente Redis com a instrumentação do Datadog
	client := redistrace.NewClient(opts, redistrace.WithServiceName(os.Getenv("DD_SERVICE")+"-redis"))

	// Executar operações Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	db_redis_con = client
	return nil
}

func NewCloseDatabaseRedisConnection() error {
	if db_redis_con == nil {
		return errors.New("redis connection is not established")
	}
	return db_redis_con.Close()
}
