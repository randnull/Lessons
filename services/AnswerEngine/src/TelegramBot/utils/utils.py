from aiogram.types import Message
from AnswerEngine.src.TelegramBot.keyboards.keyboards import student_start_keyboard, tutor_start_keyboard


async def welcome_student(message: Message) -> None:
    welcome_text = (
        f"Привет, <b>{message.from_user.full_name}</b>! 👋\n"
        "Я твой помощник в поиске репетиторов. Здесь ты можешь:\n"
        "📚 <b>Создать заказ</b> — опиши, какого репетитора ты ищешь, и я найду подходящих.\n"
        "👩‍🏫 <b>Список репетиторов</b> — посмотри доступных репетиторов и выбери лучшего.\n"
        "💬 <b>Общение</b> — общайся с репетиторами, которые готовю взять твой заказ. Никто не увидит их личные контакты — чат начну только я.\n"
        "Выбери действие ниже, чтобы начать!"
    )
    await message.answer(welcome_text, parse_mode="HTML", reply_markup=student_start_keyboard())
    print(f"Ученик {message.from_user.id} запустил бота")

async def welcome_tutor(message: Message) -> None:
    welcome_text = (
        f"Привет, <b>{message.from_user.full_name}</b>! 👋\n"
        "Я твой помощник в поиске учеников. Здесь ты можешь:\n"
        "📋 <b>Найти учеников</b> — разнообразные заказы уже ждут тебя!\n"
        "📚 <b>Список заказов</b> — смотри доступные заказы от учеников и выбирай подходящие.\n"
        "💬 <b>Общение</b> — общайся с учениками внутри телеграмм.\n"
        "Выбери действие ниже, чтобы начать!"
    )
    await message.answer(welcome_text, parse_mode="HTML", reply_markup=tutor_start_keyboard())
    print(f"Репетитор {message.from_user.id} запустил бота")
