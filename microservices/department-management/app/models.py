from sqlalchemy import Column, String, Text
from app.database import Base

class Department(Base):
    __tablename__ = "departments"

    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    description = Column(Text, nullable=True)

    def __repr__(self):
        return f"<Department {self.nombre}>"
