from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker, declarative_base

from AnswerEngine.src.config.settings import settings

DATABASE_URL = f"postgresql+asyncpg://{settings.ANSWER_DB_USER}:{settings.ANSWER_DB_PASSWORD}@{settings.ANSWER_DB_HOST}/{settings.ANSWER_DB_NAME}"

engine = create_async_engine(DATABASE_URL, echo=True, future=True)

async_session = sessionmaker(bind=engine, class_=AsyncSession)

Base = declarative_base()


async def get_session() -> AsyncSession:
    async with async_session() as session:
        yield session