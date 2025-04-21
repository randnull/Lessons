from AnswerEngine.src.TelegramBot.botStudent import bot_student
from AnswerEngine.src.TelegramBot.botTutor import bot_tutor
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto

from AnswerEngine.src.config.settings import settings


async def proceed_order(order_create: NewOrderDto) -> None:
    message = (
        f"<b>Вы создали новый заказ: {order_create.order_name}!</b>\n\n"
        "📩 <i>Мы уведомим вас, как только найдем подходящего исполнителя.</i>"
    )

    await bot_student.send_message(chat_id=str(order_create.student_id), text=message, parse_mode="html")


async def proceed_response(response: ResponseDto) -> None:
    messageStudent = (
        f"<b>На ваш заказ: {response.order_name} появился новый отклик!</b>\n\n"
    )

    messageTutor = (
        f"<b>Вы успешно откликнулись на заказ: {response.order_name}!</b>\n\n"
    )

    await bot_student.send_message(chat_id=response.student_id, text=messageStudent, parse_mode="html")
    await bot_tutor.send_message(chat_id=response.tutor_id, text=messageTutor, parse_mode="html")

async def proceed_suggest(suggest_order: SuggestDto) -> None:

    message = (
        f"<b>Вам предложили заказ: {suggest_order.order_name}</b>\n\n"
        f"<b>Описание: {suggest_order.description}</b>\n\n"
        f"<b>Цена: {suggest_order.min_price} - {suggest_order.max_price} </b>\n\n"
    )

    await bot_tutor.send_message(chat_id=suggest_order.tutor_telegram_id, text=message, parse_mode="html")


