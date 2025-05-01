from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, UUID4, Field


class OrderDto(BaseModel):
    order_id: UUID4
    order_name: str
    student_id: int
    status: str
    created_at: datetime

    @classmethod
    def to_dto(cls, OrderDao):
        return OrderDto(
            order_id=OrderDao.order_id,
            order_name=OrderDao.order_name,
            student_id=OrderDao.student_id,
            status=OrderDao.status,
            created_at=OrderDao.created_at
        )


class NewOrderDto(BaseModel):
    order_id: UUID4
    student_id: int
    order_name: str
    tags: list
    status: str

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

class SelectedDto(BaseModel):
    order_id: UUID4
    order_name: str
    student_telegram_id: int
    tutor_telegram_id: int
    response_id: UUID4

class ReviewDto(BaseModel):
    review_id: UUID4
    response_id: UUID4
    order_id: UUID4
    order_name: str
    tutor_telegram_id: int

class AddResponseDto(BaseModel):
    tutor_telegram_id: int
    response_count: int

class NewTagDto(BaseModel):
    tag_name: str

class OrderTagDto(BaseModel):
    order_id: UUID4
    tag_id: UUID4


class TutorTagDto(BaseModel):
    tutor_id: int
    tag_id: UUID4
