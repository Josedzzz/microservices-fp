from sqlalchemy import Column, String, Text
from app.database import Base
from pydantic import BaseModel
from typing import Any, List, Optional

class Department(Base):
    __tablename__ = "departments"

    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    description = Column(Text, nullable=True)

    def __repr__(self):
        return f"<Department {self.name}>"

class ErrorDetail(BaseModel):
    loc: Optional[List[str]] = None  # Location of the error (e.g., ["body", "id"])
    msg: str  # Error message
    type: Optional[str] = None  # Error type
    ctx: Optional[Any] = None  # Additional context (optional)

class ErrorResponse(BaseModel):
    code: str  # Custom error code
    message: str  # Human-readable error message
    details: Optional[List[ErrorDetail]] = None  # Validation error details (optional)

