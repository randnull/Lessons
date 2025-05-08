from aiogram import Router, F
from aiogram.filters import CommandStart, Command
from aiogram.types import Message

from AnswerEngine.src.TelegramBot.utils.utils import welcome_tutor, help_command_tutor

from AnswerEngine.src.TelegramBot.handlers.payment import send_invoice_handler, \
    success_payment_handler, pay_support_handler

from AnswerEngine.src.logger.logger import logger

tutor_router = Router()

@tutor_router.message(CommandStart())
async def cmd_start(message: Message):
    await welcome_tutor(message)


@tutor_router.message(Command("help"))
async def cmd_help(message: Message):
    logger.info(f"cmd_help run by user {message.from_user.id}")
    await help_command_tutor(message)


@tutor_router.message(Command("buy"))
async def handle_donate(message: Message):
    logger.info(f"handle_donate run by user {message.from_user.id}")
    await send_invoice_handler(message)


@tutor_router.message(Command("paysupport"))
async def pay_support(message: Message):
    logger.info(f"pay_support run by user {message.from_user.id}")
    await pay_support_handler(message)


@tutor_router.message(F.successful_payment)
async def success_payment(message: Message):
    logger.info(f"success_payment run by user {message.from_user.id}")
    await success_payment_handler(message)

@tutor_router.message(F.text.startswith("/"))
async def unknown_command(message: Message):
    logger.info(f"[tutor] unknown_command run by user {message.from_user.id} with command: {message.text}")
    await message.answer("Извините, такой команды не существует.")