package store

import (
	"context"

	"github.com/jitendravee/clean_go/internals/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Users interface {
		Create(context.Context, *models.User) (*models.User, error)
		Get(context.Context) error
	}
}

func NewStorage(db *mongo.Database) Storage {
	return Storage{
		Users: &UserStore{db},
	}
}
