package factories

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

var db_mongo_con *mongo.Client

// Abre conexão com MongoDB com suporte ao tracing do Datadog
func NewDatabaseMongoOpenConnection() error {
	// Obter URI do MongoDB do ambiente
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not defined")
		return errors.New("MONGO_URI environment variable is not defined")
	}

	// Configurar as opções de conexão do MongoDB
	clientOptions := options.Client().ApplyURI(uri)

	// Adicionar o monitor para instrumentação com Datadog
	clientOptions.Monitor = mongotrace.NewMonitor(
		mongotrace.WithServiceName(os.Getenv("DD_SERVICE") + "-mongo"),
	)

	// Conectar ao MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB: ", err)
		return err
	}

	// Testar a conexão com um ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Erro ao realizar ping no MongoDB: ", err)
		return err
	}

	db_mongo_con = client
	return nil
}

// Fecha a conexão com o MongoDB
func NewCloseDatabaseMongoConnection() error {
	if db_mongo_con == nil {
		return errors.New("MongoDB connection not established")
	}

	// Fechar a conexão com o MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db_mongo_con.Disconnect(ctx)
}
