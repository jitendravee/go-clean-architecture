package store

import (
	"context"
	"fmt"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	db *mongo.Database
}

func (s *UserStore) Create(ctx context.Context, user *models.User) (*models.User, error) {
	collection := s.db.Collection("users")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not insert user : %v", user)
	}
	return user, nil

}

func (s *UserStore) Get(ctx context.Context) error {
	return nil
}
