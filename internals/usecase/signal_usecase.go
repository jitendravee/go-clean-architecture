package usecase

import (
	"context"

	"github.com/jitendravee/clean_go/internals/models"
	"github.com/jitendravee/clean_go/internals/repository"
)

type SignalUseCase struct {
	SignalRepo repository.SignalRepo
}

func NewSignalUseCase(repo repository.SignalRepo) *SignalUseCase {
	return &SignalUseCase{
		SignalRepo: repo,
	}
}

func (uc *SignalUseCase) CreateGroupSignal(ctx context.Context, data *models.GroupSignal) (*models.GroupSignal, error) {
	return uc.SignalRepo.CreateGroupSignal(ctx, data)
}

func (uc *SignalUseCase) GetAllSignal(ctx context.Context) (*[]models.GroupSignal, error) {
	return uc.SignalRepo.GetAllSignal(ctx)
}
