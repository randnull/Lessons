from pydantic import Field
from pydantic_settings import BaseSettings
import os

from AnswerEngine.src.logger.logger import logger


class Settings(BaseSettings):
    BOT_TOKEN: str = os.getenv('BOT_STUDENT_TOKEN')
    BOT_TOKEN_TUTOR: str = os.getenv('BOT_TUTOR_TOKEN')
    FQND_HOST: str = os.getenv('FQND_HOST')
    ANSWER_DB_USER: str = os.getenv('ANSWER_DB_USER')
    ANSWER_DB_PASSWORD: str = os.getenv('ANSWER_DB_PASSWORD')
    ANSWER_DB_NAME: str = os.getenv('ANSWER_DB_NAME')
    ANSWER_DB_HOST: str = os.getenv('ANSWER_DB_HOST')
    ADMIN_USER: int = os.getenv('ADMIN_USER')
    SERVER_PORT: int = os.getenv('ANSWER_SERVER_PORT')
    GRPC_HOST: str = os.getenv('GRPCUSERHOST')
    GRPC_PORT: str = os.getenv('GRPCUSERPORT')
    MQUSER: str = os.getenv('RABBITMQ_USER')
    MQPASSWORD: str = os.getenv('RABBITMQ_PASSWORD')
    MQHOST: str = os.getenv('RABBITMQ_HOST')
    MQPORT: int = int(os.getenv('RABBITMQ_PORT'))
    DELAY_SCHEDULE: int = os.getenv('DELAY_SCHEDULE')
    SUPPORT_CHANNEL: str = os.getenv('SUPPORT_CHANNEL')

    def get_webhook_url(self) -> str:
        return f"https://{self.FQND_HOST}/webhook"

settings = Settings()
