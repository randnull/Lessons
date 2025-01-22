from aiogram.types import Message
from AnswerEngine.src.TelegramBot.keyboards.keyboards import main_keyboard

async def send_text(message: Message) -> None:
    await message.answer(
        f"<b>{message.from_user.full_name}</b>\n"
        "Чем я могу помочь вам сегодня?",
        reply_markup=main_keyboard()
    )
