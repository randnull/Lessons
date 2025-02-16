from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    BOT_TOKEN: str = os.environ['BOT_TOKEN'] # "7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"
    FQND_HOST: str = "google.com" #os.environ['FQND_HOST']
    ANSWER_DB_USER: str = os.environ['ANSWER_DB_USER'] # "postgres"
    ANSWER_DB_PASSWORD: str = os.environ['ANSWER_DB_PASSWORD'] #"postgres"
    ANSWER_DB_NAME: str = os.environ['ANSWER_DB_NAME']# "answer_engine_database" #
    ANSWER_DB_HOST: str = os.environ['ANSWER_DB_HOST'] # "127.0.0.1:5432"
    ADMIN_USER: int = 506645542 #os.environ['ADMIN_USER']

    def get_webhook_url(self) -> str:
        return f"https://{self.FQND_HOST}/webhook"

settings = Settings()