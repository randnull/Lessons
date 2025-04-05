from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    BOT_TOKEN: str = os.getenv('BOT_TOKEN', "7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0")
    FQND_HOST: str = os.getenv('FQND_HOST', "google.com")
    ANSWER_DB_USER: str = os.getenv('ANSWER_DB_USER', "postgres")
    ANSWER_DB_PASSWORD: str = os.getenv('ANSWER_DB_PASSWORD', "postgres")
    ANSWER_DB_NAME: str = os.getenv('ANSWER_DB_NAME', "response_database")
    ANSWER_DB_HOST: str = os.getenv('ANSWER_DB_HOST', "127.0.0.1:5432")
    ADMIN_USER: int = os.getenv('ADMIN_USER', "506645542")

    def get_webhook_url(self) -> str:
        return f"https://{self.FQND_HOST}/webhook"

settings = Settings()

#postgresql-answer