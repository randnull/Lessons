import json

from aio_pika import connect, IncomingMessage

from AnswerEngine.src.functions.order_functions import create_new_order, change_order_status_to_selected
from AnswerEngine.src.functions.tags_functions import update_tags
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto, TagChangeDto, SelectedDto, \
    ReviewDto, AddResponseDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_order, proceed_response, proceed_suggest, \
    proceed_order_to_tutors, proceed_selected, proceed_review, proceed_add_response


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
        selected_order = SelectedDto(**body)

        await change_order_status_to_selected(selected_order.order_id)

        await proceed_selected(selected_order)

async def review_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        new_review = ReviewDto(**body)

        await proceed_review(new_review)

async def add_response_func(message: IncomingMessage):
    async with message.process():
        body = json.loads(message.body.decode())
        add_responses = AddResponseDto(**body)

        await proceed_add_response(add_responses)
