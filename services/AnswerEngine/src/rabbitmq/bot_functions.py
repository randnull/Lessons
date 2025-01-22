from AnswerEngine.src.TelegramBot.bot import bot
from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto


async def proceed_order(order_create: OrderEngineDto):
    # parsing ...

    # request to db

    message = (
        f"🎉 <b>Ваш заказ {order_create.title} создан!</b>\n\n"
        "<b>Информация о вашей записи:</b>\n"
        f"<b>Цена: {order_create.min_price} - {order_create.max_price} </b>\n"
    )

    print(order_create)

    await bot.send_message(chat_id=order_create.student_id, text=message)
