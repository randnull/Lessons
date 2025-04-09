from aiogram.types import ReplyKeyboardMarkup, WebAppInfo, InlineKeyboardMarkup
from aiogram.utils.keyboard import ReplyKeyboardBuilder, InlineKeyboardBuilder

def student_start_keyboard() -> ReplyKeyboardMarkup:
    builder = ReplyKeyboardBuilder()
    # builder.button(text="Мои заявки")
    builder.button(
        text="Найти репетитора",
        web_app=WebAppInfo(url="https://lessonsmy.tech/")
    )
    builder.adjust(1)
    return builder.as_markup(resize_keyboard=True)

def tutor_start_keyboard() -> ReplyKeyboardMarkup:
    builder = ReplyKeyboardBuilder()
    # builder.button(text="Мои заявки")
    builder.button(
        text="Найти учеников",
        web_app=WebAppInfo(url="https://lessonsmy.tech/tutors")
    )
    builder.adjust(1)
    return builder.as_markup(resize_keyboard=True)


def payment_keyboard():
    builder = InlineKeyboardBuilder()

    builder.button(text="5 откликов", callback_data="sub_5")
    builder.button(text="10 откликов", callback_data="sub_10")
    builder.button(text="15 откликов", callback_data="sub_15")
    builder.button(text="30 откликов", callback_data="sub_30")

    builder.adjust(2)

    return builder.as_markup()


