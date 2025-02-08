from pydantic import BaseModel, Field, UUID4

class OrderEngineDto(BaseModel):
    id: UUID4 = Field(..., title="Order ID")
    student_id: int = Field(..., title="Order Title")
    title: str = Field(..., title="Order Title")
    description: str = Field(..., title="Order Description")
    min_price: int = Field(..., title="Order MinPrice")
    max_price: int = Field(..., title="Order MaxPrice")
    tags: list = Field(..., title="Order Tags")
    status: str = Field(..., title="Order Status")
    # chat_id: int = Field(..., title="Chat ID")

    #     name: str = Field(..., min_length=2, max_length=50, description="Имя клиента")
    @classmethod
    def to_dto(cls, order_engine_dao): # -> OrderEngineDto
        return OrderEngineDto(
            id=order_engine_dao.id,
            student_id=order_engine_dao.student_id,
            title=order_engine_dao.title,
            description=order_engine_dao.description,
            min_price=order_engine_dao.min_price,
            max_price=order_engine_dao.max_price,
            tags=order_engine_dao.tags,
            status=order_engine_dao.status
        )


