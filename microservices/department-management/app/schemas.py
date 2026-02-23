from pydantic import BaseModel, Field
from typing import Optional

class DepartmentBase(BaseModel):
    name: str = Field(..., description="Department name")
    description: Optional[str] = Field(None, description="Department description")

class DepartmentCreate(DepartmentBase):
    pass

class Department(DepartmentBase):
    id: str = Field(..., description="Department unique identifier")
    
    class Config:
        from_attributes = True

class DepartmentResponse(BaseModel):
    message: str
    data: Department
