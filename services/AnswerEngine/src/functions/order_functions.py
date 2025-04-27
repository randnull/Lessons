from typing import List

from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.models.dao_table.dao import OrderDao, TagDao, OrderTagDao, TutorTagDao
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

        NoExistTags = [
            NewTagDto(tag_name=tag)
            for tag in order_data.tags
            if tag not in existing_tag_names
        ]

        await tags_repository.create_many(NoExistTags)

    async with async_session() as session:
        order = OrderDto(
            order_id=order_data.order_id,
            order_name=order_data.order_name,
            student_id=order_data.student_id,
            status=order_data.status,
        )

        order_repository = Repository[OrderDao](OrderDao, session)
        tags_order_tags_repository = Repository[OrderTagDao](OrderTagDao, session)

        await order_repository.create(order)

        existing_tag = await tags_repository.get_tags_ids_by_name(order_data.tags)

        OrderTagsDtos = list()

        for tag in existing_tag:
            OrderTagsDtos.append(
                OrderTagDto(
                    order_id=order_data.order_id,
                    tag_id=tag.id,
                )
            )
            final_tags.append(tag.id)

        await tags_order_tags_repository.create_many(OrderTagsDtos)

    async with async_session() as session:
        tutor_tags_repository = Repository[TutorTagDao](TutorTagDao, session)

        tutors = await tutor_tags_repository.get_tutors_by_tags(final_tags)

        final_tutors_id = [tutor.tutor_id for tutor in tutors]

    return final_tutors_id