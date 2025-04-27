from typing import List, Optional

from pydantic import BaseModel, UUID4, Field


class OrderDto(BaseModel):
    order_id: UUID4
    order_name: str
    student_id: int
    status: str


class NewOrderDto(BaseModel):
    order_id: UUID4 = Field(..., title="Order ID")
    student_id: int = Field(..., title="Student ID")
    order_name: str = Field(..., title="Order title")
    tags: list = Field(..., title="Order Tags")
    status: str = Field(..., title="Order Status")


class ResponseDto(BaseModel):
    response_id: UUID4
    tutor_id: int
    student_id: int
    order_id: UUID4
    order_name: str

class SuggestDto(BaseModel):
    order_id: UUID4
    tutor_telegram_id: int
    order_name: str
    description: str
    min_price: int
    max_price: int

class TagChangeDto(BaseModel):
    tutor_telegram_id: int
    tags: List[str]

class TagDto(BaseModel):
    id: Optional[UUID4]
    tag_name: str

    @classmethod
    def to_dto(cls, TagDao):
        return TagDto(
            id=TagDao.id,
            tag_name=TagDao.tag_name
        )


class NewTagDto(BaseModel):
    tag_name: str


class OrderTagDto(BaseModel):
    order_id: UUID4
    tag_id: UUID4


class TutorTagDto(BaseModel):
    tutor_id: int
    tag_id: UUID4
