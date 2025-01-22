from aiogram import Router, F
from aiogram.filters import CommandStart
from aiogram.types import Message

from AnswerEngine.src.TelegramBot.utils.utils import send_text

main_router = Router()

@main_router.message(CommandStart())
async def cmd_start(message: Message):
    print(message)
    await send_text(message)
