from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.models.dao_table.dao import OrderDao, TagDao, OrderTagDao, TutorTagDao
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, OrderDto, TagDto, OrderTagDto, NewTagDto, TagChangeDto, \
    TutorTagDto


async def update_tags(new_tags: TagChangeDto) -> None:
    async with async_session() as session:
        tags_repository = Repository[TagDao](TagDao, session)

        existing_tag = await tags_repository.get_tags_ids_by_name(new_tags.tags)
        existing_tag_names = [tag.tag_name for tag in existing_tag]

        NoExistTags = [
            NewTagDto(tag_name=tag)
            for tag in new_tags.tags
            if tag not in existing_tag_names
        ]

        await tags_repository.create_many(NoExistTags)

    async with async_session() as session:
        tags_repository = Repository[TagDao](TagDao, session)
        tutor_tags_repository = Repository[TutorTagDao](TutorTagDao, session)

        updated_tags = await tags_repository.get_tags_ids_by_name(new_tags.tags)
        updated_tag_ids = {tag.id for tag in updated_tags}

        current_tags = await tutor_tags_repository.get_tags_from_tutor(new_tags.tutor_telegram_id)
        current_tag_ids = {tag.tag_id for tag in current_tags}

        tags_to_add = updated_tag_ids - current_tag_ids
        tags_to_delete = current_tag_ids - updated_tag_ids

        if tags_to_delete:
            await tutor_tags_repository.delete_many_by_conditions(
                tutor_id=new_tags.tutor_telegram_id,
                tag_ids=tags_to_delete
            )

        tutor_tag_dtos = [
            TutorTagDto(
                tutor_id=new_tags.tutor_telegram_id,
                tag_id=tag_id
            )
            for tag_id in tags_to_add
        ]

        if tutor_tag_dtos:
            await tutor_tags_repository.create_many(tutor_tag_dtos)
