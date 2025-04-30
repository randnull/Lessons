from datetime import datetime
from typing import List
from uuid import UUID

from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.logger.logger import logger
from AnswerEngine.src.models.dao_table.dao import OrderDao, TagDao, OrderTagDao, TutorTagDao, OrderStatus
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, OrderDto, TagDto, OrderTagDto, NewTagDto


async def create_new_order(order_data: NewOrderDto) -> List[int]:
    final_tags = list()
    final_tutors_id = list()

    lower_case_tags = [tag.lower() for tag in order_data.tags]
    order_data.tags = lower_case_tags

    async with async_session() as session:
        tags_repository = Repository[TagDao](TagDao, session)

        existing_tag = await tags_repository.get_tags_ids_by_name(order_data.tags)
        existing_tag_names = [tag.tag_name for tag in existing_tag]

        tags_to_add = [
            NewTagDto(tag_name=tag)
            for tag in order_data.tags
            if tag not in existing_tag_names
        ]

        if tags_to_add:
            logger.info(f"add {tags_to_add} to database")
            await tags_repository.create_many(tags_to_add)

        order = OrderDto(
            order_id=order_data.order_id,
            order_name=order_data.order_name,
            student_id=order_data.student_id,
            status=OrderStatus.NEW,
            created_at=datetime.now(),
        )

        order_repository = Repository[OrderDao](OrderDao, session)
        tags_order_tags_repository = Repository[OrderTagDao](OrderTagDao, session)

        try:
            await order_repository.create(order)
            logger.info(f"create order {order.order_id}")
        except Exception as ex:
            logger.error(f"error creating new order: {order_data.order_id}. Error: {ex}")

        existing_tag = await tags_repository.get_tags_ids_by_name(order_data.tags)

        order_tags_dtos = list()

        for tag in existing_tag:
            order_tags_dtos.append(
                OrderTagDto(
                    order_id=order_data.order_id,
                    tag_id=tag.id,
                )
            )
            final_tags.append(tag.id)

        if tags_order_tags_repository:
            logger.info(f"add tags {final_tags} to order: {order_data.order_id}")
            await tags_order_tags_repository.create_many(order_tags_dtos)

        tutor_tags_repository = Repository[TutorTagDao](TutorTagDao, session)

        tutors = await tutor_tags_repository.get_tutors_by_tags(final_tags)

        final_tutors_id = [tutor.tutor_id for tutor in tutors]

    return final_tutors_id


async def change_order_status_to_selected(orderID: UUID):
    async with async_session() as session:
        order_repository = Repository[OrderDao](OrderDao, session)

        order = await order_repository.get(orderID)

        print(order)

        orderDto = OrderDto.to_dto(order)

        print(orderDto)
        if orderDto.status != OrderStatus.NEW:
            logger.error(f"error: cannot change order to selected from status: {orderDto.status}. OrderID: {orderID}")
            return

        await order_repository.change_status(orderID)
