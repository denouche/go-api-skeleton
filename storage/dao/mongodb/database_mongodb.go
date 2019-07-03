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

const (
	mongoWriteErrorDuplicate      = 11000
	mongoWriteErrorDuplicateOther = 11001
)

type DatabaseMongoDB struct {
	client       *mongo.Client
	databaseName string
}

func handleWriteException(e mongo.WriteException) error {
	if len(e.WriteErrors) > 0 {
		switch e.WriteErrors[0].Code {
		case mongoWriteErrorDuplicate:
			fallthrough
		case mongoWriteErrorDuplicateOther:
			return dao.NewDAOError(dao.ErrTypeDuplicate, e)
		}
	}
	return e
}

func NewDatabaseMongoDB(connectionURI, dbName string) dao.Database {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Unable to get a connection to mongodb")
	}

	for {
		err = client.Ping(ctx, readpref.Primary())
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
		if err != nil {
			utils.GetLogger().WithError(err).Error("Unable to ping mongodb, waiting 2s before retrying...")
			time.Sleep(2 * time.Second)
			continue
		}
		break
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
