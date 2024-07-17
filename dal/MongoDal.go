package dal

import (
	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"
)

type MongoSettings struct {
	MongoDBAddress    string
	LoggingDatabase   string
	TxnCollection     string
	EcomTxnCollection string
	Username          string
	Password          string
	Timeout           int
}

var (
	mongoSettings MongoSettings
	mongoDBClient *mongo.Client
)

func SetMongoSettings(s MongoSettings) {
	mongoSettings = s
}

func ConnectToMongo() {
	logging.Information("Initialising Connection to Mongo")

	var err error
	mongoDBClient, err = attemptConnection()
	if err != nil {
		// Inability to establish a mongo connection constitutes a critical error
		criticalError := rpcHelp.BuildMongoCritError(err.Error())
		logging.Error(criticalError)

		return
	}

	logging.Information("Connection To Mongo Successful")
}

func GetMongoClient() (*mongo.Client, error) {
	err := pingServer()
	if err != nil {
		logging.Information("Connection to database lost, attempting to re-connect...")
		mongoDBClient, err = attemptConnection()
		if err != nil {
			// Inability to establish a mongo connection constitutes a critical error
			criticalError := rpcHelp.BuildMongoCritError(err.Error())
			logging.Error(criticalError)
		}
	}

	return mongoDBClient, err
}

func attemptConnection() (*mongo.Client, error) {
	timeout := time.Duration(mongoSettings.Timeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	clientOptions := &options.ClientOptions{
		Hosts: strings.Split(mongoSettings.MongoDBAddress, ","),
		Auth: &options.Credential{
			Username:   mongoSettings.Username,
			Password:   mongoSettings.Password,
			AuthSource: mongoSettings.LoggingDatabase,
		},
		Timeout: &timeout,
	}

	return mongo.Connect(ctx, clientOptions)
}

func pingServer() error {
	if mongoDBClient == nil {
		return errors.New("no connection to Mongo database")
	} else {
		return mongoDBClient.Ping(context.TODO(), readpref.Primary())
	}
}

func CloseMongo() {
	if mongoDBClient.Ping(context.TODO(), readpref.Primary()) == nil {
		if err := mongoDBClient.Disconnect(context.TODO()); err != nil {
			criticalError := rpcHelp.BuildMongoCritError(err.Error())
			logging.Error(criticalError)
		}
	}
}

func AggregateMongoQuery(ctx context.Context, filters []bson.M) (results []interface{}, err error) {

	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}

	cursor, err := client.Database(mongoSettings.LoggingDatabase).Collection(mongoSettings.TxnCollection).Aggregate(ctx, filters)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
