package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignalRepo interface {
	CreateGroupSignal(context.Context, *models.GroupSignal) (*models.GroupSignal, error)
	GetAllSignal(context.Context) (*models.SignalGroup, error)
	GetGroupSignalById(context.Context, string) (*models.GroupSignal, error)
	UpdateVechileCountBySignalId(context.Context, *models.UpdateVehicleCountRequest, string, string) (*models.GroupSignal, error)
}

type MongoSignalRepo struct {
	db *mongo.Database
}

func NewSignalRepo(db *mongo.Database) *MongoSignalRepo {
	return &MongoSignalRepo{db}
}
func (r *MongoSignalRepo) UpdateVechileCountBySignalId(ctx context.Context, updateCountRequest *models.UpdateVehicleCountRequest, groupId string, signalId string) (*models.GroupSignal, error) {
	collection := r.db.Collection("signals")

	objectGroupId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		log.Printf("Error converting groupId to ObjectID: %v\n", err)
		return nil, fmt.Errorf("invalid groupId format: %w", err)
	}

	filter := bson.M{
		"_id":         objectGroupId,
		"signals._id": signalId,
	}

	update := bson.M{
		"$set": bson.M{"signals.$.vehicle_count": updateCountRequest.VehicleCount},
	}

	result := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			log.Printf("No document found matching the filter.\n")
			return nil, fmt.Errorf("signal not found")
		}
		return nil, fmt.Errorf("could not update the signal count: %v", result.Err())
	}

	var updatedDoc bson.M
	if err := result.Decode(&updatedDoc); err != nil {
		log.Printf("Error decoding the updated document: %v\n", err)
		return nil, fmt.Errorf("could not decode the updated document: %v", err)
	}

	var updatedSignal models.GroupSignal
	if err := result.Decode(&updatedSignal); err != nil {
		log.Printf("Error decoding the updated signal: %v\n", err)
		return nil, fmt.Errorf("could not decode the signal: %v", err)
	}

	return &updatedSignal, nil
}

func (r *MongoSignalRepo) CreateGroupSignal(ctx context.Context, data *models.GroupSignal) (*models.GroupSignal, error) {
	collection := r.db.Collection("signals")

	for i := range data.Signals {
		if data.Signals[i].SingleSignalId == "" {
			data.Signals[i].SingleSignalId = primitive.NewObjectID().Hex()
		}
	}

	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("could not insert the data: %v", err)
	}

	data.GroupSignalId = insertResult.InsertedID.(primitive.ObjectID).Hex()

	return data, nil
}

func (r *MongoSignalRepo) GetGroupSignalById(ctx context.Context, id string) (*models.GroupSignal, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID to ObjectID: %v\n", err)
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	collection := r.db.Collection("signals")

	filter := bson.M{"_id": objectId}

	log.Printf("MongoDB filter: %+v\n", filter)

	var groupSignal models.GroupSignal

	err = collection.FindOne(ctx, filter).Decode(&groupSignal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No document found for ID: %s\n", id)
			return nil, nil
		}
		log.Printf("Error finding document for ID: %s, error: %v\n", id, err)
		return nil, err
	}

	return &groupSignal, nil
}

func (r *MongoSignalRepo) GetAllSignal(ctx context.Context) (*models.SignalGroup, error) {
	collection := r.db.Collection("signals")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not get the data: %v", err)
	}
	defer cursor.Close(ctx)

	var signalsData models.SignalGroup

	for cursor.Next(ctx) {
		var groupSignal models.GroupSignal
		if err := cursor.Decode(&groupSignal); err != nil {
			return nil, fmt.Errorf("could not decode document: %v", err)
		}
		signalsData.SignalGroup = append(signalsData.SignalGroup, groupSignal)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &signalsData, nil
}
