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
	UpdateVechileCountBySignalId(context.Context, *models.UpdateSignalCountGroup, string) (*models.GroupSignal, error)
}

type MongoSignalRepo struct {
	db *mongo.Database
}

func NewSignalRepo(db *mongo.Database) *MongoSignalRepo {
	return &MongoSignalRepo{db}
}

func (r *MongoSignalRepo) UpdateVechileCountBySignalId(ctx context.Context, updateCountRequest *models.UpdateSignalCountGroup, groupId string) (*models.GroupSignal, error) {
	collection := r.db.Collection("signals")

	// Fetch the group signals document
	var groupSignal models.GroupSignal
	objectGroupId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return nil, fmt.Errorf("invalid groupId format: %w", err)
	}

	filter := bson.M{"_id": objectGroupId}
	if err := collection.FindOne(ctx, filter).Decode(&groupSignal); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}

	// Update vehicle counts and calculate signal durations
	for _, signalUpdate := range updateCountRequest.Signals {
		for i, signal := range groupSignal.Signals {
			if signal.SingleSignalId == signalUpdate.SignalSingleId {
				groupSignal.Signals[i].VehicleCount = signalUpdate.VehicleCount
			}
		}
	}
	groupSignal.Signals = calculateSignalDurationsBasedOnCount(groupSignal.Signals)

	// Update the signals in the database
	if _, err := collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"signals": groupSignal.Signals}}); err != nil {
		return nil, fmt.Errorf("failed to update group signals: %w", err)
	}

	return &groupSignal, nil
}

// Helper function to calculate signal durations for intersections
func calculateSignalDurationsBasedOnCount(signals []models.SingleSignal) []models.SingleSignal {
	const (
		minGreen   = 10
		maxGreen   = 60
		yellowTime = 5
		defaultRed = 120
	)

	noOfSignals := len(signals)

	// Calculate green and yellow durations for each signal
	for i := range signals {
		greenTime := signals[i].VehicleCount * 2 // Dynamic green time calculation

		// Enforce minimum and maximum bounds for green time
		if greenTime < minGreen {
			greenTime = minGreen
		} else if greenTime > maxGreen {
			greenTime = maxGreen
		}

		signals[i].GreenDuration = greenTime
		signals[i].YellowDuration = yellowTime
	}

	// Calculate red durations for synchronization
	// Calculate red durations for synchronization
	for i := 0; i < noOfSignals; i++ {
		if i == 0 {
			// For the first signal, red duration is set to 0 or any default value
			signals[i].RedDuration = 0
		} else {
			// For subsequent signals, red duration is the red duration of the previous signal + the yellow and green durations of the previous signal
			if i == noOfSignals-1 {
				// For the last signal, red duration is the red duration of the previous signal + its own green duration
				signals[i].RedDuration = signals[i-1].RedDuration + signals[i-1].GreenDuration
			} else {
				// For all other signals, red duration is calculated as usual
				signals[i].RedDuration = signals[i-1].RedDuration + signals[i-1].YellowDuration + signals[i-1].GreenDuration
			}
		}
	}

	return signals
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
