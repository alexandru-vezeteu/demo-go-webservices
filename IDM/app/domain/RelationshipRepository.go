package domain

import "context"

type ObjectReference struct {
	ObjectType string
	ObjectID   string
}

type SubjectReference struct {
	SubjectType string
	SubjectID   string
	Relation    *string
}

type RelationshipTuple struct {
	Resource ObjectReference
	Relation string
	Subject  SubjectReference
}

type RelationshipFilter struct {
	ResourceFilter *ObjectReference
	RelationFilter *string
	SubjectFilter  *SubjectReference
}

type RelationshipRepository interface {
	WriteRelationship(ctx context.Context, tuple RelationshipTuple) error
	DeleteRelationship(ctx context.Context, tuple RelationshipTuple) error
	ReadRelationships(ctx context.Context, filter RelationshipFilter, limit int) ([]RelationshipTuple, error)
	CheckPermission(ctx context.Context, resource ObjectReference, permission string, subject SubjectReference) (bool, error)
}
