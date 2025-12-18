package usecase

import (
	"context"
	"idmService/domain"
)

type CheckPermissionResult struct {
	Permitted bool
	Message   string
}

type CheckPermissionUseCase interface {
	Execute(ctx context.Context, resource domain.ObjectReference, permission string, subject domain.SubjectReference) (*CheckPermissionResult, error)
}

type checkPermissionUseCase struct {
	relationshipRepo domain.RelationshipRepository
}

func NewCheckPermissionUseCase(relationshipRepo domain.RelationshipRepository) CheckPermissionUseCase {
	return &checkPermissionUseCase{
		relationshipRepo: relationshipRepo,
	}
}

func (uc *checkPermissionUseCase) Execute(ctx context.Context, resource domain.ObjectReference, permission string, subject domain.SubjectReference) (*CheckPermissionResult, error) {
	permitted, err := uc.relationshipRepo.CheckPermission(ctx, resource, permission, subject)
	if err != nil {
		return &CheckPermissionResult{
			Permitted: false,
			Message:   "Failed to check permission",
		}, err
	}

	message := "Permission denied"
	if permitted {
		message = "Permission granted"
	}

	return &CheckPermissionResult{
		Permitted: permitted,
		Message:   message,
	}, nil
}
