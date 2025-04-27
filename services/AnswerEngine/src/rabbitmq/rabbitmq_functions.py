import json

from aio_pika import connect, IncomingMessage

from AnswerEngine.src.functions.order_functions import create_new_order
from AnswerEngine.src.functions.tags_functions import update_tags
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto, TagChangeDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_order, proceed_response, proceed_suggest, proceed_order_to_tutors

async def new_order_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_order: NewOrderDto = NewOrderDto(**body)

        tutors_to_notify = await create_new_order(new_order)

        await proceed_order(new_order)
        await proceed_order_to_tutors(new_order, tutors_to_notify)

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

async def tutors_change_tags_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_suggest = TagChangeDto(**body)

        await update_tags(new_suggest)

async def selected_order_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        # new_suggest = SuggestDto(**body)

        # await update_tags(new_suggest)