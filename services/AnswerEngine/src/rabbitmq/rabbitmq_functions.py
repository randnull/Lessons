import json

from aio_pika import connect, IncomingMessage

from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto
from AnswerEngine.src.models.dto_table.response_dto import ResponsesDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_order, proceed_response
from AnswerEngine.src.functions.answer_functions import create_new_answer

async def user_func(message: IncomingMessage):
    async with message.process():
        print(message.body)

        body = json.loads(message.body.decode())
        new_order: OrderEngineDto = OrderEngineDto(**body)

        await create_new_answer(new_order)

        await proceed_order(new_order)


async def response_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_response = ResponsesDto(**body)

        await proceed_response(new_response)
