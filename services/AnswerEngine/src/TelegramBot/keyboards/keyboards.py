from aiogram.types import ReplyKeyboardMarkup, WebAppInfo, InlineKeyboardMarkup
from aiogram.utils.keyboard import ReplyKeyboardBuilder, InlineKeyboardBuilder

def main_keyboard() -> ReplyKeyboardMarkup:
    keyboard = ReplyKeyboardBuilder()
    keyboard.button(text="Some text")
    keyboard.button(text="Some text",)
    keyboard.adjust(1)
    return keyboard.as_markup(resize_keyboard=True)