from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.models.dao_table.order_engine_dao import OrderEngineDao
from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto


async def create_new_answer(order_data: OrderEngineDto):
    async with async_session() as session:
        order_engine_repository = Repository[OrderEngineDao](OrderEngineDao, session)

        response_id = await order_engine_repository.create(order_data)
