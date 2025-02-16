from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    BOT_TOKEN: str = "7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"#os.environ['BOT_TOKEN']
    FQND_HOST: str = "google.com" #os.environ['FQND_HOST']
    ANSWER_DB_USER: str = "postgres" #os.environ['ANSWER_DB_USER']
    ANSWER_DB_PASSWORD: str = "postgres" #os.environ['ANSWER_DB_PASSWORD']
    ANSWER_DB_NAME: str = "answer_engine_database" #os.environ['ANSWER_DB_NAME']
    ANSWER_DB_HOST: str = "127.0.0.1:5433" #os.environ['ANSWER_DB_HOST']
    ADMIN_USER: int = 506645542 #os.environ['ADMIN_USER']

    def get_webhook_url(self) -> str:
        return f"https://{self.FQND_HOST}/webhook"

settings = Settings()