from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    BOT_TOKEN: str = "7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"#Field(..., env="BOT_TOKEN")

    def get_webhook_url(self) -> str:
        return f"https://li9rrj-109-252-122-97.ru.tuna.am/webhook"

settings = Settings()