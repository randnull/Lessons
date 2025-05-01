from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    BOT_TOKEN: str = os.getenv('BOT_TOKEN', "7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0")
    BOT_TOKEN_TUTOR: str = os.getenv('BOT_TOKEN_TUTOR', "7236922684:AAFqILxlgUMlqmVIGmg7rkjmhnt3qyFlJ4k")
    FQND_HOST: str = os.getenv('FQND_HOST', "lessonsmy.tech")
    ANSWER_DB_USER: str = os.getenv('ANSWER_DB_USER', "postgres")
    ANSWER_DB_PASSWORD: str = os.getenv('ANSWER_DB_PASSWORD', "postgres")
    ANSWER_DB_NAME: str = os.getenv('ANSWER_DB_NAME', "response_database")
    ANSWER_DB_HOST: str = os.getenv('ANSWER_DB_HOST', "127.0.0.1:5434")
    ADMIN_USER: int = os.getenv('ADMIN_USER', "506645542")
    SERVER_PORT: int = os.getenv('ANSWER_SERVER_PORT', "7090")
    GRPC_HOST: str = os.getenv('GRPCUSERHOST', '127.0.0.1')
    GRPC_PORT: str = os.getenv('GRPCUSERPORT', '2000')
    MQUSER: str = os.getenv('RABBITMQ_USER', "guest")
    MQPASSWORD: str = os.getenv('RABBITMQ_PASSWORD', 'guest')
    MQHOST: str = os.getenv('RABBITMQ_HOST', '127.0.0.1')
    MQPORT: int = int(os.getenv('RABBITMQ_PORT', 5672))
    DELAY_SCHEDULE: int = os.getenv('DELAY_SCHEDULE', 10)

    def get_webhook_url(self) -> str:
            return f"https://{self.FQND_HOST}/webhook"

settings = Settings()
