package repository

import (
	"context"
	"fmt"
	"log"
	"math"

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

	// Convert groupId to ObjectID
	objectGroupId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		log.Printf("Error converting groupId to ObjectID: %v\n", err)
		return nil, fmt.Errorf("invalid groupId format: %w", err)
	}

	// Fetch the existing group signal document
	var groupSignal models.GroupSignal
	err = collection.FindOne(ctx, bson.M{"_id": objectGroupId}).Decode(&groupSignal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Group not found for groupId: %s\n", groupId)
			return nil, fmt.Errorf("group not found")
		}
		log.Printf("Error fetching group: %v\n", err)
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}

	// Update vehicle counts based on the request
	totalVehicleCount := 0
	for _, signalUpdate := range updateCountRequest.Signals {
		for i, signal := range groupSignal.Signals {
			if signal.SingleSignalId == signalUpdate.SignalSingleId {
				groupSignal.Signals[i].VehicleCount = signalUpdate.VehicleCount
			}
			totalVehicleCount += groupSignal.Signals[i].VehicleCount
		}
	}

	// Recalculate durations using the improved algorithm
	groupSignal.Signals = calculateSignalDurations(120, groupSignal.Signals, totalVehicleCount)

	// Update the signals in the database
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectGroupId},
		bson.M{"$set": bson.M{"signals": groupSignal.Signals}},
	)
	if err != nil {
		log.Printf("Error updating group signals: %v\n", err)
		return nil, fmt.Errorf("failed to update group signals: %w", err)
	}

	return &groupSignal, nil
}

// Helper function for duration calculation
func calculateSignalDurations(totalCycle int, signals []models.SingleSignal, totalVehicleCount int) []models.SingleSignal {
	minGreen := 10
	minYellow := 3
	maxGreen := totalCycle - minYellow - 2

	// Step 1: Calculate proportional green durations
	greenDurations := make([]float64, len(signals))
	for i, signal := range signals {
		proportionalGreen := (float64(signal.VehicleCount) / float64(totalVehicleCount)) * float64(totalCycle)
		greenDurations[i] = math.Max(float64(minGreen), math.Min(proportionalGreen, float64(maxGreen)))

		log.Printf("Signal ID: %s, Vehicle Count: %d, Proportional Green: %f, Total Vehicles: %d",
			signal.SingleSignalId, signal.VehicleCount, proportionalGreen, totalVehicleCount)
	}

	// Step 2: Adjust durations if total exceeds available time
	totalGreenTime := 0.0
	for _, green := range greenDurations {
		totalGreenTime += green
	}

	availableGreenTime := float64(totalCycle - len(signals)*minYellow)
	if totalGreenTime > availableGreenTime {
		scalingFactor := availableGreenTime / totalGreenTime
		for i := range greenDurations {
			greenDurations[i] *= scalingFactor
		}
	}

	// Step 3: Assign durations back to signals
	for i, signal := range signals {
		signal.GreenDuration = int(greenDurations[i])
		signal.YellowDuration = int(math.Max(float64(minYellow), greenDurations[i]*0.1))
		signal.RedDuration = totalCycle - signal.GreenDuration - signal.YellowDuration

		log.Printf("Signal ID: %s, Final Green Duration: %d, Yellow Duration: %d, Red Duration: %d",
			signal.SingleSignalId, signal.GreenDuration, signal.YellowDuration, signal.RedDuration)

		signals[i] = signal
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
