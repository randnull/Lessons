# from sqlalchemy import Column, Integer, String, DateTime, Float
# from common.database_connection.base import Base
from sqlalchemy import DateTime, Enum, UUID, func, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship
import enum
from AnswerEngine.common.database_connection.base import Base

import datetime


class OrderEngineDao(Base):
    __tablename__ = "order_engine"

    class OrderStatusEnum(enum.Enum):
        new: str = "New"
        proceed: str = "Proceed"
        closed: str = "Closed"

    id: Mapped[UUID] = mapped_column(UUID, primary_key=True, server_default=func.get_random_uuid())

    status: Mapped[OrderStatusEnum] = mapped_column(Enum(OrderStatusEnum), nullable=False)
    order_id: Mapped[UUID] = mapped_column(UUID, nullable=False)

    final_response_id: Mapped[UUID] = mapped_column(UUID, ForeignKey('responses.id'))

    created_at: Mapped[DateTime] = mapped_column(DateTime, nullable=False)
    updated_at: Mapped[DateTime] = mapped_column(DateTime, nullable=False)

    responses: Mapped[list["ResponsesDao"]] = relationship(back_populates="order_engine")

    @classmethod
    def to_dao(cls, order_engine_dto): #-> OrderEngineDao
        return OrderEngineDao(
            status=order_engine_dto.status,
            order_id=order_engine_dto.id,
            created_at=datetime.datetime.now(),
            updated_at=datetime.datetime.now()
        )

    # id: UUID4 = Field(..., title="Chat ID")
    # student_id: int = Field(..., title="Order Title")
    # title: str = Field(..., title="Order Title")
    # description: str = Field(..., title="Order Description")
    # min_price: int = Field(..., title="Order MinPrice")
    # max_price: int = Field(..., title="Order MaxPrice")
    # tags: list = Field(..., title="Order Tags")
    # status: str = Field(..., title="Order Status")
    # # chat_id: int = Field(..., title="Chat ID")
    #

class ResponsesDao(Base):
    __tablename__ = "responses"

    id: Mapped[UUID] = mapped_column(UUID, primary_key=True)
    order_engine_id: Mapped[UUID] = mapped_column(UUID, nullable=False)

    tutor_id: Mapped[UUID] = mapped_column(UUID, nullable=False)

    response_time: Mapped[DateTime] = mapped_column(DateTime, nullable=False)

    order_engine: Mapped["OrderEngineDao"] = relationship(back_populates="responses")
