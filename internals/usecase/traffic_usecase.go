package usecase

import (
	"context"

	"github.com/jitendravee/clean_go/internals/models"
	"github.com/jitendravee/clean_go/internals/repository"
)

type TrafficUseCase struct {
	repo repository.TrafficRepo
}

func NewTrafficUseCase(repo repository.TrafficRepo) *TrafficUseCase {
	return &TrafficUseCase{
		repo: repo,
	}
}

func (uc *TrafficUseCase) Create(ctx context.Context, trafficData *models.Traffic) (*models.Traffic, error) {
	return uc.repo.Create(ctx, trafficData)
}
