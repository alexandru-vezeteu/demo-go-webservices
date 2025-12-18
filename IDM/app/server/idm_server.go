package server

import (
	"context"

	"idmService/domain"
	pb "idmService/proto"
	"idmService/usecase"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	loginUseCase              usecase.LoginUseCase
	verifyTokenUseCase        usecase.VerifyTokenUseCase
	revokeTokenUseCase        usecase.RevokeTokenUseCase
	checkPermissionUseCase    usecase.CheckPermissionUseCase
	writeRelationshipsUseCase usecase.WriteRelationshipsUseCase
	readRelationshipsUseCase  usecase.ReadRelationshipsUseCase
}

func NewIdentityServer(
	loginUseCase usecase.LoginUseCase,
	verifyTokenUseCase usecase.VerifyTokenUseCase,
	revokeTokenUseCase usecase.RevokeTokenUseCase,
	checkPermissionUseCase usecase.CheckPermissionUseCase,
	writeRelationshipsUseCase usecase.WriteRelationshipsUseCase,
	readRelationshipsUseCase usecase.ReadRelationshipsUseCase,
) *IdentityServer {
	return &IdentityServer{
		loginUseCase:              loginUseCase,
		verifyTokenUseCase:        verifyTokenUseCase,
		revokeTokenUseCase:        revokeTokenUseCase,
		checkPermissionUseCase:    checkPermissionUseCase,
		writeRelationshipsUseCase: writeRelationshipsUseCase,
		readRelationshipsUseCase:  readRelationshipsUseCase,
	}
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.loginUseCase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: "Internal error",
		}, nil
	}

	return &pb.LoginResponse{
		Success: result.Success,
		Token:   result.Token,
		Message: result.Message,
		UserId:  result.UserID,
		Role:    result.Role,
		Email:   result.Email,
	}, nil
}

func (s *IdentityServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	result, err := s.verifyTokenUseCase.Execute(ctx, req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid:   false,
			Message: "Internal error",
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:       result.Valid,
		Email:       result.Email,
		Message:     result.Message,
		UserId:      result.UserID,
		Role:        result.Role,
		Issuer:      result.Issuer,
		ExpiresAt:   result.ExpiresAt,
		Expired:     result.Expired,
		Blacklisted: result.Blacklisted,
	}, nil
}

func (s *IdentityServer) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	result, err := s.revokeTokenUseCase.Execute(ctx, req.Token)
	if err != nil {
		return &pb.RevokeTokenResponse{
			Success: false,
			Message: "Internal error",
		}, nil
	}

	return &pb.RevokeTokenResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

func (s *IdentityServer) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	resource := domain.ObjectReference{
		ObjectType: req.Resource.ObjectType,
		ObjectID:   req.Resource.ObjectId,
	}

	subject := domain.SubjectReference{
		SubjectType: req.Subject.SubjectType,
		SubjectID:   req.Subject.SubjectId,
	}
	if req.Subject.Relation != nil {
		subject.Relation = req.Subject.Relation
	}

	result, err := s.checkPermissionUseCase.Execute(ctx, resource, req.Permission, subject)
	if err != nil {
		return &pb.CheckPermissionResponse{
			Permitted: false,
			Message:   "Internal error",
		}, nil
	}

	return &pb.CheckPermissionResponse{
		Permitted: result.Permitted,
		Message:   result.Message,
	}, nil
}

func (s *IdentityServer) WriteRelationships(ctx context.Context, req *pb.WriteRelationshipsRequest) (*pb.WriteRelationshipsResponse, error) {
	updates := make([]usecase.RelationshipUpdate, len(req.Updates))
	for i, update := range req.Updates {
		var relation *string
		if update.Relationship.Subject.Relation != nil {
			relation = update.Relationship.Subject.Relation
		}

		updates[i] = usecase.RelationshipUpdate{
			Operation: usecase.RelationshipUpdateOperation(update.Operation),
			Relationship: domain.RelationshipTuple{
				Resource: domain.ObjectReference{
					ObjectType: update.Relationship.Resource.ObjectType,
					ObjectID:   update.Relationship.Resource.ObjectId,
				},
				Relation: update.Relationship.Relation,
				Subject: domain.SubjectReference{
					SubjectType: update.Relationship.Subject.SubjectType,
					SubjectID:   update.Relationship.Subject.SubjectId,
					Relation:    relation,
				},
			},
		}
	}

	result, err := s.writeRelationshipsUseCase.Execute(ctx, updates)
	if err != nil {
		return &pb.WriteRelationshipsResponse{
			Success:              false,
			Message:              "Internal error",
			RelationshipsWritten: 0,
		}, nil
	}

	return &pb.WriteRelationshipsResponse{
		Success:              result.Success,
		Message:              result.Message,
		RelationshipsWritten: int32(result.RelationshipsWritten),
	}, nil
}

func (s *IdentityServer) ReadRelationships(ctx context.Context, req *pb.ReadRelationshipsRequest) (*pb.ReadRelationshipsResponse, error) {
	filter := domain.RelationshipFilter{}

	if req.Filter.ResourceFilter != nil {
		filter.ResourceFilter = &domain.ObjectReference{
			ObjectType: req.Filter.ResourceFilter.ObjectType,
			ObjectID:   req.Filter.ResourceFilter.ObjectId,
		}
	}

	if req.Filter.RelationFilter != nil {
		filter.RelationFilter = req.Filter.RelationFilter
	}

	if req.Filter.SubjectFilter != nil {
		var relation *string
		if req.Filter.SubjectFilter.Relation != nil {
			relation = req.Filter.SubjectFilter.Relation
		}
		filter.SubjectFilter = &domain.SubjectReference{
			SubjectType: req.Filter.SubjectFilter.SubjectType,
			SubjectID:   req.Filter.SubjectFilter.SubjectId,
			Relation:    relation,
		}
	}

	limit := 100
	if req.Limit != nil {
		limit = int(*req.Limit)
	}

	result, err := s.readRelationshipsUseCase.Execute(ctx, filter, limit)
	if err != nil {
		return &pb.ReadRelationshipsResponse{
			Relationships: []*pb.RelationshipTuple{},
			TotalCount:    0,
		}, nil
	}

	relationships := make([]*pb.RelationshipTuple, len(result.Relationships))
	for i, tuple := range result.Relationships {
		var subjectRelation *string
		if tuple.Subject.Relation != nil {
			subjectRelation = tuple.Subject.Relation
		}

		relationships[i] = &pb.RelationshipTuple{
			Resource: &pb.ObjectReference{
				ObjectType: tuple.Resource.ObjectType,
				ObjectId:   tuple.Resource.ObjectID,
			},
			Relation: tuple.Relation,
			Subject: &pb.SubjectReference{
				SubjectType: tuple.Subject.SubjectType,
				SubjectId:   tuple.Subject.SubjectID,
				Relation:    subjectRelation,
			},
		}
	}

	return &pb.ReadRelationshipsResponse{
		Relationships: relationships,
		TotalCount:    int32(result.TotalCount),
	}, nil
}
