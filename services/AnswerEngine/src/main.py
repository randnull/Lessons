import asyncio

import uvicorn

from fastapi import FastAPI
from contextlib import asynccontextmanager

from AnswerEngine.src.config.settings import settings

from AnswerEngine.src.TelegramBot.botStudent import dp_student, bot_student, start_student, stop_student
from AnswerEngine.src.TelegramBot.botTutor import dp_tutor, bot_tutor, start_tutor, stop_tutor
from AnswerEngine.src.controllers.webhook import webhook_router
from AnswerEngine.src.rabbitmq.rabbitmq_consumer import OrderConsumer, ResponseConsumer, SuggestConsumer

tags = [
    {
        "name": "Telegram Bot",
        "description": "Telegram TelegramBot"
    }
]

@asynccontextmanager
async def lifespan(app: FastAPI):
    webhook_url = settings.get_webhook_url()

    webhook_url_student = f"{webhook_url}/student"
    webhook_url_tutor = f"{webhook_url}/tutor"
    await start_student()
    await start_tutor()
    print("Student bot started.")
    print("Tutor bot started.")

    await bot_student.set_webhook(url=webhook_url_student, allowed_updates=dp_student.resolve_used_update_types(), drop_pending_updates=True)
    await bot_tutor.set_webhook(url=webhook_url_tutor, allowed_updates=dp_tutor.resolve_used_update_types(), drop_pending_updates=True)
    print("Student: ", webhook_url_student)
    print("Tutor: ", webhook_url_tutor)
    await OrderConsumer.connect()
    await ResponseConsumer.connect()
    await SuggestConsumer.connect()
    asyncio.create_task(OrderConsumer.consume())
    asyncio.create_task(ResponseConsumer.consume())
    asyncio.create_task(SuggestConsumer.consume())
    yield
    await OrderConsumer.disconnect()
    await ResponseConsumer.disconnect()
    await SuggestConsumer.disconnect()
    await stop_student()
    await stop_tutor()
    await bot_student.delete_webhook()
    await bot_tutor.delete_webhook()

app = FastAPI(openapi_tags=tags, lifespan=lifespan)

app.include_router(webhook_router)

if __name__ == "__main__":
    print("Starting AnswerEngine1")
    uvicorn.run("main:app", host="0.0.0.0", port=settings.SERVER_PORT) # , workers=5
