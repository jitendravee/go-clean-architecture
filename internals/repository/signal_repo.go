package repository

import (
	"context"
	"fmt"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignalRepo interface {
	CreateGroupSignal(context.Context, *models.GroupSignal) (*models.GroupSignal, error)
	GetAllSignal(context.Context) (*[]models.GroupSignal, error)
}

type MongoSignalRepo struct {
	db *mongo.Database
}

func NewSignalRepo(db *mongo.Database) *MongoSignalRepo {
	return &MongoSignalRepo{db}
}

func (r *MongoSignalRepo) CreateGroupSignal(ctx context.Context, data *models.GroupSignal) (*models.GroupSignal, error) {

	collection := r.db.Collection("signals")
	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("could not insert the data: %v", err)

	}
	data.GroupSignalId = insertResult.InsertedID.(primitive.ObjectID).Hex()
	return data, nil
}
func (r *MongoSignalRepo) GetAllSignal(ctx context.Context) (*[]models.GroupSignal, error) {
	return nil, nil
}
