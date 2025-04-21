from aiogram.types import ReplyKeyboardMarkup, WebAppInfo, InlineKeyboardMarkup, InlineKeyboardButton
from aiogram.utils.keyboard import ReplyKeyboardBuilder, InlineKeyboardBuilder

def student_start_keyboard() -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(
            text="🔍 Найти репетитора",
            web_app=WebAppInfo(url="https://lessonsmy.tech/")
        )]
    ])
    return keyboard

def tutor_start_keyboard() -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(
            text="🔍 Найти учеников",
            web_app=WebAppInfo(url="https://lessonsmy.tech/reps")
        )]
    ])
    return keyboard

def tutors_start_keyboard() -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(
            text="🔍 Найти учеников",
            web_app=WebAppInfo(url="https://lessonsmy.tech/reps")
        )]
    ])
    return keyboard


def payment_keyboard():
    builder = InlineKeyboardBuilder()

    builder.button(text="5 откликов", callback_data="sub_5")
    builder.button(text="10 откликов", callback_data="sub_10")
    builder.button(text="15 откликов", callback_data="sub_15")
    builder.button(text="30 откликов", callback_data="sub_30")

    builder.adjust(2)

    return builder.as_markup()


