package usecase

import (
	"context"
	"fmt"
	"idmService/domain"
)

type RelationshipUpdateOperation int

const (
	OperationCreate RelationshipUpdateOperation = 1
	OperationDelete RelationshipUpdateOperation = 2
)

type RelationshipUpdate struct {
	Operation    RelationshipUpdateOperation
	Relationship domain.RelationshipTuple
}

type WriteRelationshipsResult struct {
	Success              bool
	Message              string
	RelationshipsWritten int
}

type WriteRelationshipsUseCase interface {
	Execute(ctx context.Context, updates []RelationshipUpdate) (*WriteRelationshipsResult, error)
}

type writeRelationshipsUseCase struct {
	relationshipRepo domain.RelationshipRepository
}

func NewWriteRelationshipsUseCase(relationshipRepo domain.RelationshipRepository) WriteRelationshipsUseCase {
	return &writeRelationshipsUseCase{
		relationshipRepo: relationshipRepo,
	}
}

func (uc *writeRelationshipsUseCase) Execute(ctx context.Context, updates []RelationshipUpdate) (*WriteRelationshipsResult, error) {
	written := 0

	for _, update := range updates {
		var err error
		switch update.Operation {
		case OperationCreate:
			err = uc.relationshipRepo.WriteRelationship(ctx, update.Relationship)
		case OperationDelete:
			err = uc.relationshipRepo.DeleteRelationship(ctx, update.Relationship)
		default:
			return &WriteRelationshipsResult{
				Success:              false,
				Message:              fmt.Sprintf("Unknown operation: %d", update.Operation),
				RelationshipsWritten: written,
			}, nil
		}

		if err != nil {
			return &WriteRelationshipsResult{
				Success:              false,
				Message:              fmt.Sprintf("Failed to write relationship: %v", err),
				RelationshipsWritten: written,
			}, err
		}
		written++
	}

	return &WriteRelationshipsResult{
		Success:              true,
		Message:              fmt.Sprintf("Successfully wrote %d relationships", written),
		RelationshipsWritten: written,
	}, nil
}
