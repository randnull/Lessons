from aiogram import Router, F
from aiogram.filters import CommandStart
from aiogram.types import Message

from AnswerEngine.src.TelegramBot.utils.utils import welcome_tutor

tutor_router = Router()

@tutor_router.message(CommandStart())
async def cmd_start(message: Message):
    await welcome_tutor(message)
