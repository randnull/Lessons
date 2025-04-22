import json

from aio_pika import connect, IncomingMessage

from AnswerEngine.src.functions.order_functions import create_new_order
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_order, proceed_response, proceed_suggest


# from AnswerEngine.src.models.dto_table.response_dto import ResponsesDto
# from AnswerEngine.src.rabbitmq.bot_functions import proceed_order, proceed_response
# from AnswerEngine.src.functions.answer_functions import create_new_answer

async def new_order_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_order: NewOrderDto = NewOrderDto(**body)

        await create_new_order(new_order)

        await proceed_order(new_order)

async def response_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_response = ResponseDto(**body)

        await proceed_response(new_response)

async def suggest_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_suggest = SuggestDto(**body)

        await proceed_suggest(new_suggest)
