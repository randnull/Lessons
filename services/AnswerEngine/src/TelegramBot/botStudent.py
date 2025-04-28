import datetime
import socket

import psutil
from aiogram import Bot, Dispatcher, F
from aiogram.filters import Command


from AnswerEngine.src.TelegramBot.routers.stundents_routers import student_router
from AnswerEngine.src.config.settings import settings
from AnswerEngine.src.logger.logger import logger

bot_student = Bot(token=settings.BOT_TOKEN)
dp_student = Dispatcher()

dp_student.include_router(student_router)

async def start_student() -> None:
    try:
        await bot_student.send_message(settings.ADMIN_USER, f"Bot Started at {datetime.datetime.now()}")
        logger.info(f"Starting student bot at {datetime.datetime.now()}")
    except Exception as ex:
        logger.error(f"Student bot error start: {ex}")


async def stop_student() -> None:
    try:
        await bot_student.send_message(settings.ADMIN_USER, f"Bot Stopped at {datetime.datetime.now()}")
        logger.info(f"Stopping student bot at {datetime.datetime.now()}")
    except Exception as ex:
        logger.error(f"Student bot error stop: {ex}")

