from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    BOT_TOKEN: str = os.environ['BOT_TOKEN']
    FQND_HOST: str = os.environ['FQND_HOST']

    def get_webhook_url(self) -> str:
        return f"https://{self.FQND_HOST}/webhook"

settings = Settings()