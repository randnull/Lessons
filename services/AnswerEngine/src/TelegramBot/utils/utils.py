from aiogram.types import Message

from AnswerEngine.src.TelegramBot.keyboards.keyboards import student_start_keyboard, tutor_start_keyboard

from AnswerEngine.src.TelegramBot.utils.text import get_faq, get_welcome_text

from AnswerEngine.src.logger.logger import logger


async def welcome_student(message: Message) -> None:
    welcome_text = get_welcome_text(message.from_user.full_name, "student")
    await message.answer(welcome_text, parse_mode="HTML", reply_markup=student_start_keyboard())
    logger.info(f"student %s run start", message.from_user.id)


async def welcome_tutor(message: Message) -> None:
    welcome_text = get_welcome_text(message.from_user.full_name, "tutor")
    await message.answer(welcome_text, parse_mode="HTML", reply_markup=tutor_start_keyboard())
    logger.info(f"tutor %s run start", message.from_user.id)


async def help_command_student(message: Message, help_text = get_faq("student")) -> None:
    await message.answer(help_text, parse_mode="HTML")
    logger.info(f"student %s run help", message.from_user.id)


async def help_command_tutor(message: Message, help_text = get_faq("tutor")) -> None:
    await message.answer(help_text, parse_mode="HTML")
    logger.info(f"tutor %s run help", message.from_user.id)
