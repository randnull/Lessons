from datetime import datetime

from sqlalchemy import DateTime, Enum, UUID, func, ForeignKey, String, BIGINT, TIMESTAMP
from sqlalchemy.orm import Mapped, mapped_column, relationship
import enum
from AnswerEngine.common.database_connection.base import Base

class OrderStatus(str, Enum):
    NEW = "New"
    SELECTED = "Selected"
    CLOSED = "Closed"

class OrderDao(Base):
    __tablename__ = 'orders'

    order_id: Mapped[UUID] = mapped_column(UUID, primary_key=True)
    order_name: Mapped[String] = mapped_column(String)
    student_id: Mapped[BIGINT] = mapped_column(BIGINT)
    status: Mapped[str] = mapped_column(String(255), nullable=False)
    created_at: Mapped[TIMESTAMP] = mapped_column(TIMESTAMP, nullable=False)

    @classmethod
    def to_dao(cls, OrderDto):
        return OrderDao(
            order_id=OrderDto.order_id,
            order_name=OrderDto.order_name,
            student_id=OrderDto.student_id,
            status=OrderDto.status,
            created_at=OrderDto.created_at
        )


class ResponseDao(Base):
    __tablename__ = 'responses'

    response_id: Mapped[UUID] = mapped_column(UUID, primary_key=True)
    tutor_id: Mapped[BIGINT] = mapped_column(BIGINT)
    order_id: Mapped[UUID] = mapped_column(UUID)


class TagDao(Base):
    __tablename__ = "tags"

    id: Mapped[UUID] = mapped_column(UUID, primary_key=True, server_default=func.get_random_uuid())
    tag_name: Mapped[str] = mapped_column(String(255), nullable=False, unique=True)

    orders: Mapped[list["OrderTagDao"]] = relationship("OrderTagDao", back_populates="tag")
    tutors: Mapped[list["TutorTagDao"]] = relationship("TutorTagDao", back_populates="tag")

    @classmethod
    def to_dao(cls, TagDto):
        return TagDao(
            tag_name=TagDto.tag_name,
        )


class OrderTagDao(Base):
    __tablename__ = "order_tags"

    order_id: Mapped[UUID] = mapped_column(UUID, primary_key=True)
    tag_id: Mapped[UUID] = mapped_column(UUID, ForeignKey("tags.id"), primary_key=True)

    tag: Mapped["TagDao"] = relationship("TagDao", back_populates="orders")

    @classmethod
    def to_dao(cls, OrderTagDto):
        return OrderTagDao(
            order_id=OrderTagDto.order_id,
            tag_id=OrderTagDto.tag_id,
        )


class TutorTagDao(Base):
    __tablename__ = "tutor_tags"

    tutor_id: Mapped[BIGINT] = mapped_column(BIGINT, primary_key=True)
    tag_id: Mapped[UUID] = mapped_column(UUID, ForeignKey("tags.id"), primary_key=True)

    tag: Mapped["TagDao"] = relationship("TagDao", back_populates="tutors")

    @classmethod
    def to_dao(cls, TutorTagDto):
        return TutorTagDao(
            tutor_id=TutorTagDto.tutor_id,
            tag_id=TutorTagDto.tag_id,
        )


class SuggestDao(Base):
    __tablename__ = "suggest_table"

    tutor_id: Mapped[BIGINT] = mapped_column(BIGINT, primary_key=True)
    order_id: Mapped[UUID] = mapped_column(UUID, primary_key=True)

    @classmethod
    def to_dao(cls, SuggestDto):
        return SuggestDao(
            tutor_id=SuggestDto.tutor_telegram_id,
            order_id=SuggestDto.order_id,
        )