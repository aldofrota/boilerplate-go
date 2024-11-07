package factories

import (
	"errors"
	"os"

	"github.com/go-redis/redis"
)

var db_redis_con *redis.Client

func NewDatabaseRedisOpenConnection() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		return errors.New("REDIS_HOST or REDIS_PORT not found in environment variables")
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
		DB:   0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	db_redis_con = client
	return nil
}

func NewCloseDatabaseRedisConnection() error {
	return db_redis_con.Close()
}
