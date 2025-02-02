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

func (uc *SignalUseCase) GetAllSignal(ctx context.Context) (*models.SignalGroup, error) {
	return uc.SignalRepo.GetAllSignal(ctx)
}

func (uc *SignalUseCase) GetGroupSignalByIdUseCase(ctx context.Context, groupId string) (*models.GroupSignal, error) {
	return uc.SignalRepo.GetGroupSignalById(ctx, groupId)
}
func (uc *SignalUseCase) UpdateVechileCountBySignalIdUseCase(ctx context.Context, updateCountRequest *models.UpdateSignalCountGroup, groupId string) (*models.GroupSignal, error) {
	return uc.SignalRepo.UpdateVechileCountBySignalId(ctx, updateCountRequest, groupId)
}

func (uc *SignalUseCase) UpdateImageUrlUseCase(ctx context.Context, imageRequestList *models.ImageRequestList, groupId string) (*models.GroupSignal, error) {
	return uc.SignalRepo.UpdateImageUrl(ctx, imageRequestList, groupId)
}
