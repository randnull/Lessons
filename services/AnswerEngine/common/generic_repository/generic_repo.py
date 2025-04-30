from typing import TypeVar, Generic, List, Any, Sequence, Optional
from uuid import UUID

from sqlalchemy import select, and_, update, func, Row, RowMapping, delete
from sqlalchemy.ext.asyncio import AsyncSession

from pydantic import BaseModel

from AnswerEngine.src.models.dao_table.dao import OrderStatus

Model = TypeVar('Model')

class Repository(Generic[Model]):
    def __init__(self, model, session: AsyncSession):
        self.__model = model
        self.__session = session

    async def create(self, model: BaseModel) -> UUID:
        result_dao = self.__model.to_dao(model)
        self.__session.add(result_dao)
        await self.__session.flush()
        str_id = str(self.__model.__table__.columns.keys()[0])
        id = getattr(result_dao, str_id)
        await self.__session.commit()
        return id

    async def create_many(self, models: List[BaseModel]):
        daos = [self.__model.to_dao(model) for model in models]
        self.__session.add_all(daos)
        await self.__session.flush()
        await self.__session.commit()

    async def get_tags_ids_by_name(self, tags: List[str]):
        resp = await self.__session.execute(select(self.__model).where(self.__model.tag_name.in_(tags)))
        return resp.scalars().all()

    async def get_tags_from_tutor(self, tutor_id):
        result = await self.__session.execute(select(self.__model).where(self.__model.tutor_id == tutor_id))
        return result.scalars().all()

    async def get_tutors_by_tags(self, tags: List[UUID]):
        resp = await self.__session.execute(select(self.__model).where(self.__model.tag_id.in_(tags)))
        return resp.scalars().all()

    async def delete_many_by_conditions(self, tutor_id, tag_ids) -> None:
         await self.__session.execute(delete(self.__model).where(
            self.__model.tutor_id == tutor_id,
            self.__model.tag_id.in_(tag_ids)
         ))
         await self.__session.commit()

    async def change_status(self, order_id: UUID) -> bool:
         await self.__session.execute(
            update(self.__model)
            .where(self.__model.order_id == order_id)
            .values(status=OrderStatus.SELECTED)
         )
         await self.__session.commit()

    async def get(self, order_id: UUID) -> Optional[Model]:
        result = await self.__session.execute(
            select(self.__model).where(self.__model.order_id == order_id)
        )
        return result.scalars().first()