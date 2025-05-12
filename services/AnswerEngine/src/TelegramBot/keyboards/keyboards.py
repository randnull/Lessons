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

def suggest_keyboard(order_id) -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(
            text="Перейти к заказу",
            web_app=WebAppInfo(url=f"https://lessonsmy.tech/reps/#/order/{order_id}")
        )]
    ])
    return keyboard

def review_keyboard(review_id) -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(
            text="Перейти к отклику",
            web_app=WebAppInfo(url=f"https://lessonsmy.tech/reps/#/response/id/{review_id}")
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

