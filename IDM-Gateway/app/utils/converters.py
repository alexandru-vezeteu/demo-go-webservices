from app.proto import idm_pb2
from app.models.authz_models import (
    ObjectReference, SubjectReference, RelationshipTuple,
    RelationshipUpdate, RelationshipOperation
)

def pydantic_object_to_proto(obj: ObjectReference) -> idm_pb2.ObjectReference:
    return idm_pb2.ObjectReference(
        object_type=obj.object_type,
        object_id=obj.object_id
    )

def pydantic_subject_to_proto(subj: SubjectReference) -> idm_pb2.SubjectReference:
    proto_subj = idm_pb2.SubjectReference(
        subject_type=subj.subject_type,
        subject_id=subj.subject_id
    )
    if subj.relation is not None:
        proto_subj.relation = subj.relation
    return proto_subj

def proto_tuple_to_pydantic(tuple_proto: idm_pb2.RelationshipTuple) -> RelationshipTuple:
    return RelationshipTuple(
        resource=ObjectReference(
            object_type=tuple_proto.resource.object_type,
            object_id=tuple_proto.resource.object_id
        ),
        relation=tuple_proto.relation,
        subject=SubjectReference(
            subject_type=tuple_proto.subject.subject_type,
            subject_id=tuple_proto.subject.subject_id,
            relation=tuple_proto.subject.relation if tuple_proto.subject.HasField("relation") else None
        )
    )

def pydantic_update_to_proto(update: RelationshipUpdate) -> idm_pb2.RelationshipUpdate:
    operation = (
        idm_pb2.OPERATION_CREATE if update.operation == RelationshipOperation.CREATE
        else idm_pb2.OPERATION_DELETE
    )
    return idm_pb2.RelationshipUpdate(
        operation=operation,
        relationship=idm_pb2.RelationshipTuple(
            resource=pydantic_object_to_proto(update.relationship.resource),
            relation=update.relationship.relation,
            subject=pydantic_subject_to_proto(update.relationship.subject)
        )
    )
