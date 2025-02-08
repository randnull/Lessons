from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession

from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import get_session


def get_repository(model):
    def get(session: AsyncSession = Depends(get_session)):
        return Repository(model, session)

    return get