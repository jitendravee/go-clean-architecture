package repository

import (
	"context"
	"fmt"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrafficRepo interface {
	Create(ctx context.Context, trafficData *models.Traffic) (*models.Traffic, error)
}

type MongoTrafficRepo struct {
	db *mongo.Database
}

func NewMongoTrafficRepo(db *mongo.Database) *MongoTrafficRepo {
	return &MongoTrafficRepo{db}
}

func (r *MongoTrafficRepo) Create(ctx context.Context, trafficData *models.Traffic) (*models.Traffic, error) {
	collection := r.db.Collection("traffic")

	_, err := collection.InsertOne(ctx, trafficData)
	if err != nil {
		return nil, fmt.Errorf("could not inset the data : %v", err)
	}
	return trafficData, nil
}
