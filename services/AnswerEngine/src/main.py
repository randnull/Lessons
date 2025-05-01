import asyncio

import uvicorn

from fastapi import FastAPI
from contextlib import asynccontextmanager

from AnswerEngine.src.config.settings import settings

from AnswerEngine.src.TelegramBot.botStudent import dp_student, bot_student, start_student, stop_student
from AnswerEngine.src.TelegramBot.botTutor import dp_tutor, bot_tutor, start_tutor, stop_tutor
from AnswerEngine.src.controllers.webhook import webhook_router
from AnswerEngine.src.rabbitmq.rabbitmq_consumer import OrderConsumer, ResponseConsumer, SuggestConsumer, \
    TagsChangeConsumer, SelectedConsumer, ReviewConsumer

from AnswerEngine.src.jobs.job import start_scheduler, stop_scheduler


from AnswerEngine.src.logger.logger import logger

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
    logger.info("student bot init")
    await start_tutor()
    logger.info("tutors bot init")
    try:
        await bot_student.set_webhook(url=webhook_url_student, allowed_updates=dp_student.resolve_used_update_types(),
                                      drop_pending_updates=True)
    except Exception as ex:
        logger.error("student bot webhook failed: %s", str(ex))
        logger.info("student bot webhook ok: %s", webhook_url_student)

    try:
        await bot_tutor.set_webhook(url=webhook_url_tutor, allowed_updates=dp_tutor.resolve_used_update_types(),
                                drop_pending_updates=True)
        logger.info("tutors bot webhook ok: %s", webhook_url_tutor)
    except Exception as ex:
        logger.error("tutor bot webhook ok: %s", str(ex))
    await OrderConsumer.connect()
    await ResponseConsumer.connect()
    await SuggestConsumer.connect()
    await TagsChangeConsumer.connect()
    await SelectedConsumer.connect()
    await ReviewConsumer.connect()

    await start_scheduler()

    asyncio.create_task(OrderConsumer.consume())
    asyncio.create_task(ResponseConsumer.consume())
    asyncio.create_task(SuggestConsumer.consume())
    asyncio.create_task(TagsChangeConsumer.consume())
    asyncio.create_task(SelectedConsumer.consume())
    asyncio.create_task(ReviewConsumer.consume())
    yield
    await OrderConsumer.disconnect()
    await ResponseConsumer.disconnect()
    await SuggestConsumer.disconnect()
    await TagsChangeConsumer.disconnect()
    await SelectedConsumer.disconnect()
    await ReviewConsumer.disconnect()
    await stop_student()
    await stop_tutor()
    await bot_student.delete_webhook()
    await bot_tutor.delete_webhook()

    await stop_scheduler()
app = FastAPI(openapi_tags=tags, lifespan=lifespan)

app.include_router(webhook_router)

if __name__ == "__main__":
    logger.info(f"Starting AnswerEngine server. Port: {settings.SERVER_PORT}")
    uvicorn.run("main:app", host="0.0.0.0", port=settings.SERVER_PORT) # , workers=5
