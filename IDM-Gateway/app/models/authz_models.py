from pydantic import BaseModel, Field
from typing import Optional, List
from enum import Enum

class ObjectReference(BaseModel):
    object_type: str = Field(..., min_length=1)
    object_id: str = Field(..., min_length=1)

class SubjectReference(BaseModel):
    subject_type: str = Field(..., min_length=1)
    subject_id: str = Field(..., min_length=1)
    relation: Optional[str] = None

class RelationshipTuple(BaseModel):
    resource: ObjectReference
    relation: str = Field(..., min_length=1)
    subject: SubjectReference

class CheckPermissionRequest(BaseModel):
    resource: ObjectReference
    permission: str = Field(..., min_length=1)
    subject: SubjectReference

class CheckPermissionResponse(BaseModel):
    permitted: bool
    message: str

class RelationshipOperation(str, Enum):
    CREATE = "CREATE"
    DELETE = "DELETE"

class RelationshipUpdate(BaseModel):
    operation: RelationshipOperation
    relationship: RelationshipTuple

class WriteRelationshipsRequest(BaseModel):
    updates: List[RelationshipUpdate] = Field(..., min_items=1)

class WriteRelationshipsResponse(BaseModel):
    success: bool
    message: str
    relationships_written: int

class ReadRelationshipsResponse(BaseModel):
    relationships: List[RelationshipTuple]
    total_count: int
