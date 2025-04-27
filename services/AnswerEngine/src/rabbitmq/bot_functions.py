from typing import List

from AnswerEngine.src.TelegramBot.botStudent import bot_student
from AnswerEngine.src.TelegramBot.botTutor import bot_tutor
from AnswerEngine.src.TelegramBot.keyboards.keyboards import suggest_keyboard
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto, TagChangeDto

from AnswerEngine.src.config.settings import settings

async def proceed_order(order_create: NewOrderDto) -> None:
    message = (
        f"<b>Вы успешно создали новый заказ: {order_create.order_name}!</b>\n\n"
        "📩 <i>Мы сообщим вам, как только подберем подходящего исполнителя.</i>"
    )

    await bot_student.send_message(chat_id=str(order_create.student_id), text=message, parse_mode="html")

async def proceed_order_to_tutors(order_create: NewOrderDto, tutors_id: List[int]) -> None:
    message = (
        f"<b>Заказ подходит вашим тегам!: {order_create.order_name}!</b>\n\n"
    )

    for tutor_id in tutors_id:
        await bot_tutor.send_message(chat_id=str(tutor_id), text=message, parse_mode="html", reply_markup=suggest_keyboard(order_create.order_id))


async def proceed_response(response: ResponseDto) -> None:
    messageStudent = (
        f"<b>У вашего заказа \"{response.order_name}\" появился новый отклик!</b>\n\n"
        "👀 <i>Вы можете рассмотреть отклик и связаться с исполнителем.</i>"
    )

    messageTutor = (
        f"<b>Вы откликнулись на заказ \"{response.order_name}\"!</b>\n\n"
        "✅ <i>Ожидайте ответа от заказчика.</i>"
    )

    await bot_student.send_message(chat_id=response.student_id, text=messageStudent, parse_mode="html")
    await bot_tutor.send_message(chat_id=response.tutor_id, text=messageTutor, parse_mode="html")

async def proceed_suggest(suggest_order: SuggestDto) -> None:
    message = (
        f"<b>Новый заказ для вас: {suggest_order.order_name}</b>\n\n"
        f"👀 <b>Описание:</b> {suggest_order.description}\n\n"
        f"<b>Бюджет:</b> {suggest_order.min_price} - {suggest_order.max_price}\n\n"
        "⚡ <i>Вы можете просмотреть заказ, нажав на кнопку ниже.</i>"
    )

    tutor_id = suggest_order.tutor_telegram_id

    await bot_tutor.send_message(chat_id=tutor_id, text=message, parse_mode="html", reply_markup=suggest_keyboard(suggest_order.order_id))
