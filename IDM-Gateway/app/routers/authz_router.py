from fastapi import APIRouter, HTTPException, Query
from typing import Optional
from app.models.authz_models import (
    CheckPermissionRequest, CheckPermissionResponse,
    WriteRelationshipsRequest, WriteRelationshipsResponse,
    ReadRelationshipsResponse
)
from app.grpc_client import grpc_client
from app.proto import idm_pb2
from app.utils.error_handler import handle_grpc_error
from app.utils.converters import (
    pydantic_object_to_proto, pydantic_subject_to_proto,
    pydantic_update_to_proto, proto_tuple_to_pydantic
)
import grpc

router = APIRouter(prefix="/api/idm/authz", tags=["Authorization"])

@router.post("/check", response_model=CheckPermissionResponse)
async def check_permission(request: CheckPermissionRequest):
    try:
        stub = grpc_client.get_stub()
        grpc_request = idm_pb2.CheckPermissionRequest(
            resource=pydantic_object_to_proto(request.resource),
            permission=request.permission,
            subject=pydantic_subject_to_proto(request.subject)
        )
        response = await stub.CheckPermission(grpc_request)

        return CheckPermissionResponse(
            permitted=response.permitted,
            message=response.message
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)

@router.post("/relationships", response_model=WriteRelationshipsResponse)
async def write_relationships(request: WriteRelationshipsRequest):
    try:
        stub = grpc_client.get_stub()
        updates = [pydantic_update_to_proto(update) for update in request.updates]
        grpc_request = idm_pb2.WriteRelationshipsRequest(updates=updates)
        response = await stub.WriteRelationships(grpc_request)

        return WriteRelationshipsResponse(
            success=response.success,
            message=response.message,
            relationships_written=response.relationships_written
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)

@router.get("/relationships", response_model=ReadRelationshipsResponse)
async def read_relationships(
    resource_type: Optional[str] = Query(None, description="Filter by resource type"),
    resource_id: Optional[str] = Query(None, description="Filter by resource ID"),
    relation: Optional[str] = Query(None, description="Filter by relation"),
    subject_type: Optional[str] = Query(None, description="Filter by subject type"),
    subject_id: Optional[str] = Query(None, description="Filter by subject ID"),
    limit: int = Query(100, ge=1, le=1000, description="Maximum number of results")
):
    try:
        stub = grpc_client.get_stub()

        filter_proto = idm_pb2.RelationshipFilter()
        if resource_type and resource_id:
            filter_proto.resource_filter.CopyFrom(
                idm_pb2.ObjectReference(object_type=resource_type, object_id=resource_id)
            )
        if relation:
            filter_proto.relation_filter = relation
        if subject_type and subject_id:
            filter_proto.subject_filter.CopyFrom(
                idm_pb2.SubjectReference(subject_type=subject_type, subject_id=subject_id)
            )

        grpc_request = idm_pb2.ReadRelationshipsRequest(
            filter=filter_proto,
            limit=limit
        )
        response = await stub.ReadRelationships(grpc_request)

        relationships = [proto_tuple_to_pydantic(t) for t in response.relationships]

        return ReadRelationshipsResponse(
            relationships=relationships,
            total_count=response.total_count
        )
    except grpc.RpcError as e:
        raise handle_grpc_error(e)
