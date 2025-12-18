package usecase

import (
	"context"
	"idmService/domain"
)

type ReadRelationshipsResult struct {
	Relationships []domain.RelationshipTuple
	TotalCount    int
}

type ReadRelationshipsUseCase interface {
	Execute(ctx context.Context, filter domain.RelationshipFilter, limit int) (*ReadRelationshipsResult, error)
}

type readRelationshipsUseCase struct {
	relationshipRepo domain.RelationshipRepository
}

func NewReadRelationshipsUseCase(relationshipRepo domain.RelationshipRepository) ReadRelationshipsUseCase {
	return &readRelationshipsUseCase{
		relationshipRepo: relationshipRepo,
	}
}

func (uc *readRelationshipsUseCase) Execute(ctx context.Context, filter domain.RelationshipFilter, limit int) (*ReadRelationshipsResult, error) {
	if limit <= 0 {
		limit = 100
	}

	relationships, err := uc.relationshipRepo.ReadRelationships(ctx, filter, limit)
	if err != nil {
		return &ReadRelationshipsResult{
			Relationships: []domain.RelationshipTuple{},
			TotalCount:    0,
		}, err
	}

	return &ReadRelationshipsResult{
		Relationships: relationships,
		TotalCount:    len(relationships),
	}, nil
}
