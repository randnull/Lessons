from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.models.dao_table.order_engine_dao import ResponsesDao
from AnswerEngine.src.models.dto_table.response_dto import ResponsesDto


async def create_new_answer(response_data: ResponsesDto):
    async with async_session() as session:
        response_repository = Repository[ResponsesDao](ResponsesDao, session)

        answer_id = await response_repository.create(response_data)
