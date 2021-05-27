package persistence

import (
	// "database/sql"
	"context"
	"fmt"
	"time"
	"todoapi/domain/model"
	"todoapi/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName string = "items"
	queryTimeout          = 5 * time.Second
)

type itemMongoRepository struct {
	collection *mongo.Collection
}

func NewItemMongoRepository(db *mongo.Database) repository.ItemRepository {
	collection := db.Collection(collectionName)
	return &itemMongoRepository{
		collection: collection,
	}
}

func (repo itemMongoRepository) Save(item model.Item) (model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	res, err := repo.collection.InsertOne(ctx, bson.D{{"name", item.Name}, {"done", item.Done}})
	item.Id = res.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		return model.Item{}, err
	}
	return item, err
}

// https://github.com/motiv-labs/janus/blob/master/pkg/plugin/basic/mongodb_repository.go パクった
func (repo itemMongoRepository) GetAll() ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.Item
	for cursor.Next(ctx) {
		var item struct {
			Id   primitive.ObjectID `bson:"_id"`
			Name string             `bson:"name"`
			Done bool               `bson:"done"`
		}
		err := cursor.Decode(&item)
		fmt.Println(err)
		if err != nil {
			return nil, err
		}
		fmt.Println(item)
		items = append(items, model.Item{
			Id:   model.ID(item.Id),
			Name: item.Name,
			Done: item.Done,
		})
	}
	return items, nil
}

func (repo itemMongoRepository) GetById(id model.ID) (model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var item model.Item
	objId, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return item, err
	}
	err = repo.collection.FindOne(ctx, bson.D{{"_id", objId}}).Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (repo itemMongoRepository) DeleteById(id model.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}
	_, err = repo.collection.DeleteOne(ctx, bson.D{{"_id", objId}})
	return err
}

func (repo itemMongoRepository) UpdateDone(id model.ID, done bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}
	_, err = repo.collection.UpdateByID(ctx, objId, bson.D{{"$set", bson.D{{"done", done}}}})
	return err
}

func (repo itemMongoRepository) DeleteDone() error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := repo.collection.DeleteMany(ctx, bson.D{{"done", true}})
	return err
}
