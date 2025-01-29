from AnswerEngine.src.TelegramBot.bot import bot
from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto
from AnswerEngine.src.models.dto_table.response_dto import ResponsesDto

from AnswerEngine.src.config.settings import settings


async def proceed_order(order_create: OrderEngineDto) -> None:
    # parsing ...

    # request to db

    message = (
        f"🎉 <b>Ваш заказ успешно создан!</b>\n\n"
        f"📌 <b>Название:</b> {order_create.title}\n"
        f"📝 <b>Описание:</b> {order_create.description}\n"
        f"💰 <b>Цена:</b> {order_create.min_price} - {order_create.max_price} ₽\n\n"
        "📩 <i>Мы уведомим вас, как только найдем подходящего исполнителя.</i>"
    )

    await bot.send_message(chat_id=order_create.student_id, text=message, parse_mode="html")


async def proceed_response(response: ResponsesDto) -> None:
    message = (
        f"Новый отклик на заказ # {response.order_id}"
    )

    print('Я на proceed response', response)
    # await bot.send_message(chat_id=response.order_id, text=message, parse_mode="html")

    await bot.send_message(chat_id=settings.ADMIN_USER, text=message, parse_mode="html")
