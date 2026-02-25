from pydantic import BaseModel, Field, field_validator
from typing import Optional

class DepartmentBase(BaseModel):
    id: str = Field(..., description="Department unique identifier (must be uppercase)")
    name: str = Field(..., description="Department name")
    description: Optional[str] = Field(None, description="Department description")

    @field_validator('id')
    @classmethod
    def id_must_be_uppercase(cls, v):
        if not v or not v.strip():
            raise ValueError("Department ID cannot be empty")
        if not v.isupper():
            raise ValueError("Department ID must be in uppercase")
        return v.strip()

    @field_validator('name')
    @classmethod
    def name_not_empty(cls, v):
        if not v or not v.strip():
            raise ValueError('Department name cannot be empty')
        return v.strip()

class DepartmentCreate(DepartmentBase):
    pass

class Department(DepartmentBase):
    id: str = Field(..., description="Department unique identifier")
    
    class Config:
        from_attributes = True

class DepartmentResponse(BaseModel):
    message: str
    data: Department

class ErrorResponse(BaseModel):
    detail: str
