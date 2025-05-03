from aiogram import Router, F
from aiogram.filters import CommandStart, Command
from aiogram.types import Message

from AnswerEngine.src.TelegramBot.utils.utils import welcome_student

student_router = Router()

@student_router.message(CommandStart())
async def cmd_start(message: Message):
    await welcome_student(message)
