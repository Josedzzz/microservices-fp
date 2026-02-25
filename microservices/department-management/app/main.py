from fastapi import FastAPI, Depends, HTTPException, status, Request
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session
from typing import List
from contextlib import asynccontextmanager

from app import crud, models, schemas
from app.database import SessionLocal, engine

# Create database tables
models.Base.metadata.create_all(bind=engine)

@asynccontextmanager
async def lifespan(_: FastAPI):
    # Startup
    print("Starting departments service...")
    yield
    # Shutdown
    print("Shutting down departments service...")

# Initialize FastAPI app
app = FastAPI(
    title="Departments Service",
    description="Microservice for managing departments",
    version="1.0.0",
    lifespan=lifespan
)

# Dependency to get DB session
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# Custom exception handler for HTTPException
@app.exception_handler(HTTPException)
async def http_exception_handler(_: Request, exc: HTTPException):
    error_response = models.ErrorResponse(
        code=str(exc.status_code),
        message=exc.detail,
        details=[
            models.ErrorDetail(
                loc=["request", "path"],
                msg=exc.detail,
                type="http_exception"
            )
        ]
    )
    return JSONResponse(
        status_code=exc.status_code,
        content=error_response.model_dump()
    )

# Custom exception handler for validation errors
from fastapi.exceptions import RequestValidationError
@app.exception_handler(RequestValidationError)
async def validation_exception_handler(_: Request, exc: RequestValidationError):
    error_details = [
        models.ErrorDetail(
            loc=error.get("loc", []),
            msg=error.get("msg", ""),
            type=error.get("type", "")
        )
        for error in exc.errors()
    ]
    error_response = models.ErrorResponse(
        code="422",
        message="Validation error",
        details=error_details
    )
    return JSONResponse(
        status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
        content=error_response.model_dump()
    )

@app.get("/departments-service/api/")
async def root():
    return {
        "service": "Departments Service",
        "version": "1.0.0",
        "endpoints": [
            "POST /departments - Create a new department",
            "GET /departments/{id} - Get department by ID",
            "GET /departments - List all departments"
        ]
    }

@app.post(
    "/departments-service/api/departments",
    response_model=schemas.DepartmentResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new department",
    description="Creates a new department with the provided information"
)
async def create_department(
    department: schemas.DepartmentCreate,
    db: Session = Depends(get_db)
):
    db_department = crud.create_department(db=db, department=department)
    return {
        "message": "Department created successfully",
        "data": db_department
    }

@app.get(
    "/departments-service/api/departments/{department_id}",
    response_model=schemas.DepartmentResponse,
    summary="Get department by ID",
    description="Retrieves a specific department by its ID"
)
async def get_department(
    department_id: str,
    db: Session = Depends(get_db)
):
    department = crud.get_department(db, department_id)
    if not department:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Department with id {department_id} not found"
        )
    return {
        "message": "Department retrieved successfully",
        "data": department
    }

@app.get(
    "/departments-service/api/departments",
    response_model=List[schemas.Department],
    summary="List all departments",
    description="Retrieves a list of all departments"
)
async def get_departments(
    skip: int = 0,
    limit: int = 100,
    db: Session = Depends(get_db)
):
    departments = crud.get_departments(db, skip=skip, limit=limit)
    return departments
