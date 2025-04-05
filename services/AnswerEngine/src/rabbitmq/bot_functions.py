from AnswerEngine.src.TelegramBot.bot import bot
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto
# from AnswerEngine.src.models.dto_table.response_dto import ResponsesDto

from AnswerEngine.src.config.settings import settings


async def proceed_order(order_create: NewOrderDto) -> None:
    message = (
        f"<b>Вы создали новый заказ: {order_create.order_name}!</b>\n\n"
        "📩 <i>Мы уведомим вас, как только найдем подходящего исполнителя.</i>"
    )

    await bot.send_message(chat_id=str(order_create.student_id), text=message, parse_mode="html")


async def proceed_response(response: ResponseDto) -> None:
    messageStudent = (
        f"<b>На ваш заказ: {response.order_name} появился новый отклик!</b>\n\n"
    )

    messageTutor = (
        f"<b>Вы успешно откликнулись на заказ: {response.order_name}!</b>\n\n"
    )

    await bot.send_message(chat_id=response.student_id, text=messageStudent, parse_mode="html")
    await bot.send_message(chat_id=response.tutor_id, text=messageTutor, parse_mode="html")
