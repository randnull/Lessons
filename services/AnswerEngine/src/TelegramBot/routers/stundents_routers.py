from aiogram import Router, F
from aiogram.filters import CommandStart, Command
from aiogram.types import Message

from AnswerEngine.src.TelegramBot.utils.utils import welcome_student, help_command_student
from AnswerEngine.src.logger.logger import logger

student_router = Router()

@student_router.message(CommandStart())
async def cmd_start(message: Message):
    await welcome_student(message)


@student_router.message(Command("help"))
async def cmd_help(message: Message):
    logger.info(f"cmd_help run by user {message.from_user.id}")
    await help_command_student(message)


@student_router.message(F.text.startswith("/"))
async def unknown_command(message: Message):
    logger.info(f"[student] unknown command run by user {message.from_user.id} with command: {message.text}")
    await message.answer("Извините, такой команды не существует.")