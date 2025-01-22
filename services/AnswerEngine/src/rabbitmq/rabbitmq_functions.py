import json

from aio_pika import connect, IncomingMessage

from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_order


async def user_func(message: IncomingMessage):
    async with message.process():
        print(message.body)

        body = json.loads(message.body.decode())
        new_order: OrderEngineDto = OrderEngineDto(**body)


        await proceed_order(new_order)