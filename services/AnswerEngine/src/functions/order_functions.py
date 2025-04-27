from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.models.dao_table.dao import OrderDao, TagDao, OrderTagDao
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, OrderDto, TagDto, OrderTagDto, NewTagDto


async def create_new_order(order_data: NewOrderDto) -> None:
    async with async_session() as session:
        tags_repository = Repository[TagDao](TagDao, session)

        existing_tag = await tags_repository.get_tags_ids_by_name(order_data.tags)
        existing_tag_names = [tag.tag_name for tag in existing_tag]

        OrderTagsDtos = [
            OrderTagDto(
                order_id=order_data.order_id,
                tag_id=tag.id,
            )
            for tag in existing_tag
        ]

        NoExistTags = [
            NewTagDto(tag_name=tag)
            for tag in order_data.tags
            if tag not in existing_tag_names
        ]

        await tags_repository.create_many(NoExistTags)

    async with async_session() as session:
        print('trying to create new answer')
        order = OrderDto(
            order_id=order_data.order_id,
            order_name=order_data.order_name,
            student_id=order_data.student_id,
            status=order_data.status,
        )

        order_repository = Repository[OrderDao](OrderDao, session)
        tags_order_tags_repository = Repository[OrderTagDao](OrderTagDao, session)

        await order_repository.create(order)
        await tags_order_tags_repository.create_many(OrderTagsDtos)
