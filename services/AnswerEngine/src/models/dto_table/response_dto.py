from pydantic import BaseModel, Field, UUID4

class ResponsesDto(BaseModel):
    # tutor_id: int = Field(..., title="Tutor ID")
    order_id: str = Field(..., title="Order ID")
