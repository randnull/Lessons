from datetime import datetime

from AnswerEngine.common.generic_repository.generic_repo import Repository
from AnswerEngine.common.database_connection.base import async_session
from AnswerEngine.src.logger.logger import logger
from AnswerEngine.src.models.dao_table.dao import OrderDao, OrderStatus
from AnswerEngine.src.models.dto_table.dto import OrderDto
from AnswerEngine.src.rabbitmq.bot_functions import proceed_need_review


async def selected_status_check():
    time_now = datetime.now()

    logger.info(f"Job status check started at {time_now}")

    async with async_session() as session:
        order_repository = Repository[OrderDao](OrderDao, session)

        orders = await order_repository.get_all_selected()

        logger.info(f"Selected {len(orders)} orders Selected.")

        for order in orders:
            order_dto: OrderDto = OrderDto.to_dto(order)

            time_diff = time_now - order_dto.created_at

            logger.debug(f"Order diff: {time_diff}")
            if time_diff.seconds >= 10: # testing
                await proceed_need_review(order_dto)
                await order_repository.change_status(order_dto.order_id, OrderStatus.CLOSED)
