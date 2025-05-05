from datetime import datetime

from aiogram import Router, F
from aiogram.filters import CommandStart, Command
from aiogram.types import Message, PreCheckoutQuery, CallbackQuery

from AnswerEngine.src.TelegramBot.utils.utils import welcome_tutor

from AnswerEngine.src.TelegramBot.handlers.payment import send_invoice_handler, pre_checkout_handler, \
    success_payment_handler, process_subscription_callback, pay_support_handler

from AnswerEngine.src.logger.logger import logger

tutor_router = Router()

@tutor_router.message(CommandStart())
async def cmd_start(message: Message):
    await welcome_tutor(message)

@tutor_router.message(Command("buy"))
async def handle_donate(message: Message):
    logger.info(f"handle_donate run at {datetime.now()} by user {message.from_user.id}")
    await send_invoice_handler(message)

@tutor_router.message(Command("paysupport"))
async def pay_support(message: Message):
    logger.info(f"pay_support run at {datetime.now()} by user {message.from_user.id}")
    await pay_support_handler(message)

@tutor_router.pre_checkout_query()
async def handle_pre_checkout(pre_checkout_query: PreCheckoutQuery):
    await pre_checkout_handler(pre_checkout_query)

@tutor_router.message(F.successful_payment)
async def handle_successful_payment(message: Message):
    await success_payment_handler(message)

@tutor_router.callback_query()
async def handle_subscription_callback(callback: CallbackQuery):
    await process_subscription_callback(callback)
