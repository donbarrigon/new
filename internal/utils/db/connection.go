package db

import (
	"context"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/logs"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Col(col string) *mongo.Collection {
	return DB.Collection(col)
}

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.DbConnectionString).SetServerAPIOptions(serverAPI)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetRetryWrites(true)
	clientOptions.SetTimeout(30 * time.Second)

	var err error

	Client, err = mongo.Connect(clientOptions)
	if err != nil {
		logs.Info("🔴💥 Fail to connect db %s: %s", config.DbName, config.DbConnectionString)
		return err
	}
	DB = Client.Database(config.DbName)

	logs.Info("🍃 Successful connection to %s: %s", config.DbName, config.DbConnectionString)
	return nil
}

func CloseMongoDB() error {
	if Client == nil {
		return nil
	}
	return Client.Disconnect(context.TODO())
}
