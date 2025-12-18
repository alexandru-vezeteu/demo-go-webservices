package authorization

import (
	"context"
	"fmt"
	"idmService/domain"
	"strings"
	"sync"
)

type InMemoryRelationshipRepository struct {
	mu            sync.RWMutex
	relationships map[string]domain.RelationshipTuple
}

func NewInMemoryRelationshipRepository() *InMemoryRelationshipRepository {
	return &InMemoryRelationshipRepository{
		relationships: make(map[string]domain.RelationshipTuple),
	}
}

func tupleKey(tuple domain.RelationshipTuple) string {
	subjectRel := ""
	if tuple.Subject.Relation != nil {
		subjectRel = "#" + *tuple.Subject.Relation
	}
	return fmt.Sprintf("%s:%s#%s@%s:%s%s",
		tuple.Resource.ObjectType,
		tuple.Resource.ObjectID,
		tuple.Relation,
		tuple.Subject.SubjectType,
		tuple.Subject.SubjectID,
		subjectRel,
	)
}

func (r *InMemoryRelationshipRepository) WriteRelationship(ctx context.Context, tuple domain.RelationshipTuple) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := tupleKey(tuple)
	r.relationships[key] = tuple
	return nil
}

func (r *InMemoryRelationshipRepository) DeleteRelationship(ctx context.Context, tuple domain.RelationshipTuple) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := tupleKey(tuple)
	delete(r.relationships, key)
	return nil
}

func (r *InMemoryRelationshipRepository) ReadRelationships(ctx context.Context, filter domain.RelationshipFilter, limit int) ([]domain.RelationshipTuple, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []domain.RelationshipTuple
	count := 0

	for _, tuple := range r.relationships {
		if matchesFilter(tuple, filter) {
			results = append(results, tuple)
			count++
			if limit > 0 && count >= limit {
				break
			}
		}
	}

	return results, nil
}

func matchesFilter(tuple domain.RelationshipTuple, filter domain.RelationshipFilter) bool {
	if filter.ResourceFilter != nil {
		if filter.ResourceFilter.ObjectType != "" && tuple.Resource.ObjectType != filter.ResourceFilter.ObjectType {
			return false
		}
		if filter.ResourceFilter.ObjectID != "" && tuple.Resource.ObjectID != filter.ResourceFilter.ObjectID {
			return false
		}
	}

	if filter.RelationFilter != nil && *filter.RelationFilter != "" {
		if tuple.Relation != *filter.RelationFilter {
			return false
		}
	}

	if filter.SubjectFilter != nil {
		if filter.SubjectFilter.SubjectType != "" && tuple.Subject.SubjectType != filter.SubjectFilter.SubjectType {
			return false
		}
		if filter.SubjectFilter.SubjectID != "" && tuple.Subject.SubjectID != filter.SubjectFilter.SubjectID {
			return false
		}
	}

	return true
}

func (r *InMemoryRelationshipRepository) CheckPermission(ctx context.Context, resource domain.ObjectReference, permission string, subject domain.SubjectReference) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	visited := make(map[string]bool)
	return r.checkPermissionRecursive(resource, permission, subject, visited), nil
}

func (r *InMemoryRelationshipRepository) checkPermissionRecursive(resource domain.ObjectReference, permission string, subject domain.SubjectReference, visited map[string]bool) bool {
	visitKey := fmt.Sprintf("%s:%s#%s@%s:%s", resource.ObjectType, resource.ObjectID, permission, subject.SubjectType, subject.SubjectID)
	if visited[visitKey] {
		return false
	}
	visited[visitKey] = true

	directRelation := strings.TrimSuffix(permission, "_permission")
	if directRelation == "" {
		directRelation = permission
	}

	for _, tuple := range r.relationships {
		if tuple.Resource.ObjectType == resource.ObjectType &&
			tuple.Resource.ObjectID == resource.ObjectID &&
			tuple.Relation == directRelation {

			if tuple.Subject.SubjectType == subject.SubjectType &&
				tuple.Subject.SubjectID == subject.SubjectID &&
				tuple.Subject.Relation == nil {
				return true
			}

			if tuple.Subject.Relation != nil {
				subResource := domain.ObjectReference{
					ObjectType: tuple.Subject.SubjectType,
					ObjectID:   tuple.Subject.SubjectID,
				}
				subSubject := domain.SubjectReference{
					SubjectType: subject.SubjectType,
					SubjectID:   subject.SubjectID,
				}
				if r.checkPermissionRecursive(subResource, *tuple.Subject.Relation, subSubject, visited) {
					return true
				}
			}
		}
	}

	return false
}
