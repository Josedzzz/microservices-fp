from fastapi import FastAPI, Depends, HTTPException, status, Request
from fastapi.responses import JSONResponse
from fastapi.openapi.docs import get_swagger_ui_html
from sqlalchemy.orm import Session
from typing import List
from contextlib import asynccontextmanager
from prometheus_fastapi_instrumentator import Instrumentator

from app import crud, models, schemas, auth
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
    openapi_url=None,
    docs_url=None,
    lifespan=lifespan,
    swagger_ui_parameters={"persistAuthorization": True},
    openapi_tags=[{"name": "Departments", "description": "Operations with departments"}],
)

# Initialize Prometheus instrumentation
Instrumentator().instrument(app).expose(app)

@app.get("/docs", include_in_schema=False)
async def custom_swagger_ui_html():
    from fastapi.openapi.docs import get_swagger_ui_html
    html = get_swagger_ui_html(
        openapi_url="./openapi.json",
        title=app.title + " - Swagger UI",
        oauth2_redirect_url=app.swagger_ui_oauth2_redirect_url,
        swagger_js_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js",
        swagger_css_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css",
    )
    # Force the URL to be relative in the HTML content
    body = html.body.decode("utf-8").replace("url: '/openapi.json'", "url: './openapi.json'")
    from fastapi.responses import HTMLResponse
    return HTMLResponse(content=body)

@app.get("/openapi.json", include_in_schema=False)
async def get_open_api_endpoint():
    from fastapi.openapi.utils import get_openapi
    return get_openapi(
        title=app.title, 
        version=app.version, 
        routes=app.routes,
        servers=[{"url": "/departments-service", "description": "API Gateway"}]
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

@app.get("/api/")
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
    "/api/departments",
    response_model=schemas.DepartmentResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new department",
    description="Creates a new department with the provided information",
    dependencies=[Depends(auth.security)]
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
    "/api/departments/{department_id}",
    response_model=schemas.DepartmentResponse,
    summary="Get department by ID",
    description="Retrieves a specific department by its ID",
    dependencies=[Depends(auth.security)]
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
    "/api/departments",
    response_model=List[schemas.Department],
    summary="List all departments",
    description="Retrieves a list of all departments",
    dependencies=[Depends(auth.security)]
)
async def get_departments(
    skip: int = 0,
    limit: int = 100,
    db: Session = Depends(get_db)
):
    departments = crud.get_departments(db, skip=skip, limit=limit)
    return departments
