from aiogram import Bot, Dispatcher

from AnswerEngine.src.config.settings import settings

bot = Bot(token=settings.BOT_TOKEN)
dp = Dispatcher()


async def start() -> None:
    try:
        await bot.send_message(506645542, "start")
    except Exception as error:
        pass
        # TODO logger


async def stop() -> None:
    try:
        await bot.send_message(506645542, "stop")
    except Exception as error:
        pass
        # TODO logger
