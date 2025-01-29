import datetime
import socket

import psutil
from aiogram import Bot, Dispatcher

from AnswerEngine.src.config.settings import settings

bot = Bot(token=settings.BOT_TOKEN)
dp = Dispatcher()


async def start() -> None:
    try:
        await bot.send_message(settings.ADMIN_USER, f"Bot Started at {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger


async def stop() -> None:
    try:
        await bot.send_message(settings.ADMIN_USER, f"BOT STOPPED. Time: {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger
