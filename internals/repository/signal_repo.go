package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignalRepo interface {
	CreateGroupSignal(context.Context, *models.GroupSignal) (*models.GroupSignal, error)
	GetAllSignal(context.Context) (*models.SignalGroup, error)
	GetGroupSignalById(context.Context, string) (*models.GroupSignal, error)
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
func (r *MongoSignalRepo) GetGroupSignalById(ctx context.Context, id string) (*models.GroupSignal, error) {
	// Log the start of the function and the ID being searched for
	log.Printf("GetGroupSignalById called with ID: %s\n", id)

	// Reference the collection in MongoDB
	collection := r.db.Collection("signals")

	// Create the filter for the query based on the ID
	filter := bson.M{"group_id": id}

	// Log the filter to ensure it's correct
	log.Printf("MongoDB filter: %+v\n", filter)

	// Declare a variable to hold the result
	var groupSignal models.GroupSignal

	// Perform the query to find the document
	err := collection.FindOne(ctx, filter).Decode(&groupSignal)

	// Check for errors after querying
	if err != nil {
		// Log the error if it occurs
		if err == mongo.ErrNoDocuments {
			log.Printf("No document found for ID: %s\n", id)
			return nil, nil
		}
		log.Printf("Error finding document for ID: %s, error: %v\n", id, err)
		return nil, err
	}

	// Log the successful retrieval of the document
	log.Printf("Found document for ID: %s: %+v\n", id, groupSignal)

	// Return the found document
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
