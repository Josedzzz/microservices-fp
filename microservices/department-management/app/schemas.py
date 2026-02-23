from pydantic import BaseModel, Field, field_validator
from typing import Optional

class DepartmentBase(BaseModel):
    name: str = Field(..., description="Department name")
    description: Optional[str] = Field(None, description="Department description")

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
