package mongodb

import (
	"time"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DatabaseMongoDB) GetAllUsers() ([]*model.User, error) {
	ctx := db.getCtx()
	cur, err := db.getSession().Collection("user").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]*model.User, 0)
	for cur.Next(ctx) {
		var result *model.User
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *DatabaseMongoDB) GetUserByID(id string) (*model.User, error) {
	ctx := db.getCtx()
	var result *model.User
	err := db.getSession().Collection("user").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DatabaseMongoDB) CreateUser(user *model.User) error {
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()

	ctx := db.getCtx()
	_, err := db.getSession().Collection("user").InsertOne(ctx, user)
	return err
}

func (db *DatabaseMongoDB) DeleteUser(id string) error {
	ctx := db.getCtx()
	_, err := db.getSession().Collection("user").DeleteOne(ctx, bson.M{"_id": id})
	if err == mongo.ErrNoDocuments {
		return dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return err
}

func (db *DatabaseMongoDB) UpdateUser(user *model.User) error {
	now := time.Now()
	user.UpdatedAt = &now

	ctx := db.getCtx()
	r, err := db.getSession().Collection("user").ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if r.MatchedCount == 0 {
		return dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return err
}
