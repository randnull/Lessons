import asyncio

import uvicorn

from fastapi import FastAPI
from contextlib import asynccontextmanager

from AnswerEngine.src.config.settings import settings

from AnswerEngine.src.TelegramBot.routers.main_router import main_router
from AnswerEngine.src.TelegramBot.bot import dp, bot, start, stop
from AnswerEngine.src.controllers.webhook import webhook_router
from AnswerEngine.src.controllers.tags_controller import tags_router
from AnswerEngine.src.rabbitmq.rabbitmq_consumer import OrderConsumer, ResponseConsumer

tags = [
    {
        "name": "Telegram Bot",
        "description": "Telegram TelegramBot"
    }
]

@asynccontextmanager
async def lifespan(app: FastAPI):
    dp.include_router(main_router)
    webhook_url = settings.get_webhook_url()
    await start()
    await bot.set_webhook(url=webhook_url, allowed_updates=dp.resolve_used_update_types(), drop_pending_updates=True)
    await OrderConsumer.connect()
    await ResponseConsumer.connect()
    asyncio.create_task(OrderConsumer.consume())
    asyncio.create_task(ResponseConsumer.consume())
    yield
    await OrderConsumer.disconnect()
    await ResponseConsumer.disconnect()
    await stop()
    await bot.delete_webhook()

app = FastAPI(openapi_tags=tags, lifespan=lifespan)

app.include_router(webhook_router)
app.include_router(tags_router)

if __name__ == "__main__":
    print("Starting AnswerEngine")
    uvicorn.run("main:app", host="0.0.0.0", port=7090) # , workers=5
