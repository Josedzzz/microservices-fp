from fastapi import FastAPI, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List
from contextlib import asynccontextmanager

from app import crud, models, schemas
from app.database import SessionLocal, engine

# Create database tables
models.Base.metadata.create_all(bind=engine)

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting departments service...")
    print(app)
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

@app.get("/departments-service/api/")
async def root():
    return {
        "service": "Departments Service",
        "version": "1.0.0",
        "endpoints": [
            "POST /departamentos - Create a new department",
            "GET /departamentos/{id} - Get department by ID",
            "GET /departamentos - List all departments"
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
    """
    Create a new department:
    
    - **name**: Department name (required)
    - **description**: Department description (optional)
    """
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
    """
    Get a department by its unique ID
    
    - **department_id**: The ID of the department to retrieve
    """
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
    """
    Get all departments with pagination:
    
    - **skip**: Number of records to skip (default: 0)
    - **limit**: Maximum number of records to return (default: 100)
    """
    departments = crud.get_departments(db, skip=skip, limit=limit)
    return departments
