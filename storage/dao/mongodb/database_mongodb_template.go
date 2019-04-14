package mongodb

import (
	"time"

	"github.com/denouche/go-api-skeleton/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"github.com/denouche/go-api-skeleton/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionTemplateName = "template"
)

func (db *DatabaseMongoDB) populateTemplateIndexes() {
	ctx := db.getCtx()
	_, err := db.getSession().Collection(collectionTemplateName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "name", Value: bsonx.Int32(1)}},
		Options: &options.IndexOptions{
			Unique: utils.NewBool(true),
		},
	})
	if err != nil {
		utils.GetLogger().WithError(err).Error("error while creating mongodb index")
	}
}

func (db *DatabaseMongoDB) GetAllTemplates() ([]*model.Template, error) {
	ctx := db.getCtx()
	cur, err := db.getSession().Collection(collectionTemplateName).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]*model.Template, 0)
	for cur.Next(ctx) {
		var result *model.Template
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

func (db *DatabaseMongoDB) GetTemplateByID(id string) (*model.Template, error) {
	ctx := db.getCtx()
	var result *model.Template
	err := db.getSession().Collection(collectionTemplateName).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DatabaseMongoDB) CreateTemplate(template *model.Template) error {
	template.ID = primitive.NewObjectID().Hex()
	template.CreatedAt = time.Now()

	ctx := db.getCtx()
	_, err := db.getSession().Collection(collectionTemplateName).InsertOne(ctx, template)
	return err
}

func (db *DatabaseMongoDB) DeleteTemplate(id string) error {
	ctx := db.getCtx()
	_, err := db.getSession().Collection(collectionTemplateName).DeleteOne(ctx, bson.M{"_id": id})
	if err == mongo.ErrNoDocuments {
		return dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return err
}

func (db *DatabaseMongoDB) UpdateTemplate(template *model.Template) error {
	now := time.Now()
	template.UpdatedAt = &now

	ctx := db.getCtx()
	r, err := db.getSession().Collection(collectionTemplateName).ReplaceOne(ctx, bson.M{"_id": template.ID}, template)
	if r.MatchedCount == 0 {
		return dao.NewDAOError(dao.ErrTypeNotFound, err)
	}
	return err
}
