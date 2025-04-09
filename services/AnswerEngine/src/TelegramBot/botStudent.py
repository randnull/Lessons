import datetime
import socket

import psutil
from aiogram import Bot, Dispatcher, F
from aiogram.filters import Command


from AnswerEngine.src.TelegramBot.routers.stundents_routers import student_router
from AnswerEngine.src.config.settings import settings

bot_student = Bot(token=settings.BOT_TOKEN)
dp_student = Dispatcher()

dp_student.include_router(student_router)

async def start_student() -> None:
    try:
        await bot_student.send_message(settings.ADMIN_USER, f"Bot Started at {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger


async def stop_student() -> None:
    try:
        await bot_student.send_message(settings.ADMIN_USER, f"BOT STOPPED. Time: {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger
