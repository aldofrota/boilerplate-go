package factories

import (
	"boilerplate-go/app/data/usecases"
	mongoHelper "boilerplate-go/app/infra/db/mongo/helpers"
	redisHelper "boilerplate-go/app/infra/db/redis/helpers"
	"boilerplate-go/app/presentation/controllers"
	"boilerplate-go/app/presentation/protocols"
)

// HealthCheck godoc
// @Summary      Validate if service is healthy
// @Description  Validate if mongo database, redis database is connected
// @Tags         Health Check
// @Accept       json
// @Produce      json
// @Success      200  {object}  protocols.HttpResponse "OK"
// @Failure      404  {object}  protocols.HttpResponse "Not Found"
// @Failure      500  {object}  protocols.HttpResponse "Internal Server Error"
// @Router       /health [get]
func NewHealthCheckControllerFactory() protocols.Controller {
	redisDatabaseIsConnectedHelper := redisHelper.NewRedisDatabaseIsConnectedHelper(db_redis_con)
	mongoDatabaseIsConnectedHelper := mongoHelper.NewMongoDatabaseIsConnectedHelper(db_mongo_con)
	validateIfHealthyService := usecases.NewValidateIfHealthyService(
		redisDatabaseIsConnectedHelper,
		mongoDatabaseIsConnectedHelper,
	)
	return controllers.NewHealthCheckController(validateIfHealthyService)
}
