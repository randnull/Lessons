from typing import TypeVar, Generic, List, Any, Sequence
from uuid import UUID

from sqlalchemy import select, and_, update, func, Row, RowMapping
from sqlalchemy.ext.asyncio import AsyncSession

from pydantic import BaseModel

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

        # stmt = select(self.__model).where(TagDao.tag_name.in_(order_data.tags))
        # result = await session.execute(stmt)
        # existing_tags = result.scalars().all()

    # async def get_by_code(self, code: str) -> Model:
    #     resp = await self.__session.execute(select(self.__model).where(self.__model.code == code))
    #     return resp.scalars().one_or_none()
    #
    # async def get_by_value(self, value_name: str, value):
    #     resp = await self.__session.execute(select(self.__model).where(getattr(self.__model, value_name) == value))
    #     return resp.scalars().all()
    #
    # async def get_one_by_value(self, value_name: str, value) -> Model:
    #     resp = await self.__session.execute(select(self.__model).where(getattr(self.__model, value_name) == value))
    #     return resp.scalars().one_or_none()
    #
    # async def get_by_values(self, dict_vals: dict) -> Model:
    #     resp = await self.__session.execute(select(self.__model).where(
    #         and_(getattr(self.__model, value_name) == val for value_name, val in dict_vals.items())))
    #     return resp.scalars().all()
    #
    # async def get_by_id(self, id: int):
    #     str_id = str(self.__model.__table__.columns.keys()[0])
    #     resp = await self.__session.execute(select(self.__model).where(getattr(self.__model, str_id) == id))
    #     return resp.scalars().one_or_none()
    #
    # async def get(self):
    #     resp = await self.__session.execute(select(self.__model))
    #     return resp.scalars().all()
    #
    # async def update_value_not_id(self, value_find_name: str, value_find, value_name: str, new_value):
    #     await self.__session.execute(
    #         update(self.__model).where(getattr(self.__model, value_find_name) == value_find).values(status=new_value))
    #     await self.__session.commit()
    #
    # async def update_value(self, id: int, value_name: str, new_value):
    #     str_id = str(self.__model.__table__.columns.keys()[0])
    #     await self.__session.execute(
    #         update(self.__model).where(getattr(self.__model, str_id) == id).values(status=new_value))
    #     await self.__session.commit()
    #
    # async def update_val(self, id: int, value_name: str, new_value):
    #     str_id = str(self.__model.__table__.columns.keys()[0])
    #     await self.__session.execute(
    #         update(self.__model).where(getattr(self.__model, str_id) == id).values({value_name: new_value}))
    #     await self.__session.commit()
    #
    # async def add(self, model: BaseModel) -> int:
    #     result_dao = self.__model.to_dao(model)
    #     self.__session.add(result_dao)
    #     await self.__session.flush()
    #     str_id = str(self.__model.__table__.columns.keys()[0])
    #     id = getattr(result_dao, str_id)
    #     await self.__session.commit()
    #     return id
    #
    # async def check(self, model: BaseModel) -> bool:
    #     req = await self.__session.execute(select(self.__model).where(self.__model.code == model.code))
    #     is_find = req.scalars().one_or_none()
    #
    #     if is_find:
    #         return True
    #     return False
    #
    # async def check_by_columns(self, another):
    #     columns = self.__model.__table__.columns.keys()
    #     id_name = str(columns[0])
    #
    #     req = await self.__session.execute(select(self.__model).where(
    #         and_(getattr(self.__model, feat) == getattr(another, feat) for feat in columns if feat != id_name)))
    #     return req.scalars().one_or_none()
    #
    # async def delete(self, code: str):
    #     req = await self.__session.execute(select(self.__model).where(self.__model.code == code))
    #     result_req = req.scalar()
    #     if result_req is not None:
    #         await self.__session.delete(result_req)
    #         await self.__session.commit()
    #
    # async def update_status(self, id: int, new_value):
    #     str_id = str(self.__model.__table__.columns.keys()[0])
    #     await self.__session.execute(
    #         update(self.__model).where(getattr(self.__model, str_id) == id).values(status=new_value))
    #     await self.__session.commit()