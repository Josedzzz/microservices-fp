from sqlalchemy.orm import Session
from fastapi import HTTPException
from app import models, schemas

def get_department(db: Session, department_id: str):
    return db.query(models.Department).filter(models.Department.id == department_id).first()

def get_departments(db: Session, skip: int = 0, limit: int = 100):
    return db.query(models.Department).offset(skip).limit(limit).all()

def create_department(db: Session, department: schemas.DepartmentCreate):
    # Check if the department already exists
    existing_department = db.query(models.Department).filter(models.Department.id == department.id).first()
    if existing_department:
        raise HTTPException(
            status_code=400,
            detail=f"Department with id '{department.id}' already exists."
        )

    db_department = models.Department(**department.model_dump())
    db.add(db_department)
    db.commit()
    db.refresh(db_department)
    return db_department
