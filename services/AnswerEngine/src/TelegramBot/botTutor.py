import datetime

from aiogram import Bot, Dispatcher

from AnswerEngine.src.TelegramBot.handlers.payment import pre_checkout_handler, process_subscription_callback
from AnswerEngine.src.TelegramBot.routers.tutors_routers import tutor_router

from AnswerEngine.src.config.settings import settings
from AnswerEngine.src.logger.logger import logger

bot_tutor = Bot(token=settings.BOT_TOKEN_TUTOR)
dp_tutor = Dispatcher()

dp_tutor.include_router(tutor_router)

dp_tutor.pre_checkout_query.register(pre_checkout_handler)
dp_tutor.callback_query.register(process_subscription_callback)


async def start_tutor() -> None:
    try:
        await bot_tutor.send_message(settings.ADMIN_USER, f"Bot Started at {datetime.datetime.now()}")
        logger.info(f"Starting tutor bot at {datetime.datetime.now()}")
    except Exception as ex:
        logger.error(f"Tutor bot error start: {ex}")

async def stop_tutor() -> None:
    try:
        await bot_tutor.send_message(settings.ADMIN_USER, f"BOT Stopped. Time: {datetime.datetime.now()}")
        logger.info(f"Stopping tutor bot at {datetime.datetime.now()}")
    except Exception as ex:
        logger.error(f"Tutor bot error stop: {ex}")
