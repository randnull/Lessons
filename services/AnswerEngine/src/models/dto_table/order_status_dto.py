from pydantic import BaseModel, Field, UUID4

class OrderEngineDto(BaseModel):
    id: UUID4 = Field(..., title="Chat ID")
    student_id: int = Field(..., title="Order Title")
    title: str = Field(..., title="Order Title")
    description: str = Field(..., title="Order Description")
    min_price: int = Field(..., title="Order MinPrice")
    max_price: int = Field(..., title="Order MaxPrice")
    tags: list = Field(..., title="Order Tags")
    status: str = Field(..., title="Order Status")
    chat_id: int = Field(..., title="Chat ID")

    pass #     name: str = Field(..., min_length=2, max_length=50, description="Имя клиента")

