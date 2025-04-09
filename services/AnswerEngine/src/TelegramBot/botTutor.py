import datetime
import socket

import psutil
from aiogram import Bot, Dispatcher, F
from aiogram.filters import Command

from AnswerEngine.src.TelegramBot.handlers.payment import send_invoice_handler, pre_checkout_handler, \
    success_payment_handler, process_subscription_callback
from AnswerEngine.src.TelegramBot.routers.tutors_routers import tutor_router

from AnswerEngine.src.config.settings import settings

bot_tutor = Bot(token=settings.BOT_TOKEN_TUTOR)
dp_tutor = Dispatcher()

dp_tutor.include_router(tutor_router)
dp_tutor.message.register(send_invoice_handler, Command(commands="donate"))
dp_tutor.pre_checkout_query.register(pre_checkout_handler)
dp_tutor.message.register(success_payment_handler, F.successful_payment)
# dp_tutor.message.register(pay_support_handler, Command(commands="paysupport"))
dp_tutor.callback_query.register(process_subscription_callback)

async def start_tutor() -> None:
    try:
        await bot_tutor.send_message(settings.ADMIN_USER, f"Bot Started at {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger


async def stop_tutor() -> None:
    try:
        await bot_tutor.send_message(settings.ADMIN_USER, f"BOT STOPPED. Time: {datetime.datetime.now()}\n Current Host: {socket.gethostname()}\n CPU: {psutil.cpu_percent()}\n Memory: {psutil.virtual_memory().percent}")
    except Exception as error:
        pass
        # TODO logger
