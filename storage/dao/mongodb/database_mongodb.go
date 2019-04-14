package mongodb

import (
	"context"
	"time"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseMongoDB struct {
	client       *mongo.Client
	databaseName string
}

func NewDatabaseMongoDB(connectionURI, dbName string) dao.Database {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Unable to get a connection to mongodb")
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Unable to ping mongodb")
	}

	result := &DatabaseMongoDB{
		client:       client,
		databaseName: dbName,
	}

	result.populateTemplateIndexes() // Template index

	return result
}

func (db *DatabaseMongoDB) getSession() *mongo.Database {
	return db.client.Database(db.databaseName)
}
func (db *DatabaseMongoDB) getCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}
