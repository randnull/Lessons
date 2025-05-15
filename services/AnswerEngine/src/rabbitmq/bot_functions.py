from typing import List

from AnswerEngine.src.TelegramBot.botStudent import bot_student
from AnswerEngine.src.TelegramBot.botTutor import bot_tutor
from AnswerEngine.src.TelegramBot.keyboards.keyboards import suggest_keyboard
from AnswerEngine.src.logger.logger import logger
from AnswerEngine.src.models.dao_table.dao import OrderStatus
from AnswerEngine.src.models.dto_table.dto import NewOrderDto, ResponseDto, SuggestDto, SelectedDto, \
    ReviewDto, OrderDto, AddResponseDto

from AnswerEngine.src.config.settings import settings


async def proceed_order(order_create: NewOrderDto) -> None:
    messageWaiting = (
        f"<b>Вы успешно создали новый заказ: {order_create.order_name}!</b>\n\n"
        "📩 <i>Заказ находится на рассмотрении у администрации</i>"
    )

    messageNew = (
        f"<b>🎉 Заказ: {order_create.order_name} был одобрен администрацией для размещения!\n\n</b>"
        "📩 <i>Мы сообщим вам, как только исполнитель отправит предложение.</i>"
    )

    message = messageWaiting if order_create.status == OrderStatus.WAITING else messageNew

    try:
        await bot_student.send_message(chat_id=str(order_create.student_id), text=message, parse_mode="html")
        logger.info(f"[NOTIFY-STUDENT] order {order_create.order_id} create to user {order_create.student_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-STUDENT] order {order_create.order_id} create to user {order_create.student_id} failed!. Error: {ex}")


async def proceed_order_to_tutors(order_create: NewOrderDto, tutors_id: List[int]) -> None:
    message = (
        f"<b>Появился новый заказ, подходящий вашим тегам: {order_create.order_name}!</b>\n\n"
    )

    for tutor_id in tutors_id:
        try:
            await bot_tutor.send_message(chat_id=str(tutor_id), text=message, parse_mode="html", reply_markup=suggest_keyboard(order_create.order_id))
            logger.info(f"[NOTIFY-TUTOR] order {order_create.order_id} create to user {tutor_id} send!")
        except Exception as ex:
            logger.error(f"[NOTIFY-TUTOR] order {order_create.order_id} create to user {tutor_id} failed!. Error: {ex}")


async def proceed_response(response: ResponseDto) -> None:
    messageStudent = (
        f"<b>У вашего заказа \"{response.order_name}\" появился новый отклик!</b>\n\n"
        "👀 <i>Вы можете рассмотреть отклик и связаться с исполнителем.</i>\n\n"
        f"⚠️ <i>В случае, если репетитор просит оплату вперед занятия, сообщите об этом в поддержку: {settings.SUPPORT_CHANNEL}!</i>"
        f"⚠️ <i>Мы настоятельно не рекомендуем производить оплату вперед занятий!</i>"
    )

    messageTutor = (
        f"<b>Вы откликнулись на заказ \"{response.order_name}\"!</b>\n\n"
        "✅ <i>Ожидайте ответа от заказчика.</i>"
    )

    try:
        await bot_student.send_message(chat_id=response.student_id, text=messageStudent, parse_mode="html")
        logger.info(f"[NOTIFY-STUDENT] response: {response.response_id} to user: {response.student_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-STUDENT] response: {response.response_id} to user: {response.student_id} failed!. Error: {ex}")

    try:
        await bot_tutor.send_message(chat_id=response.tutor_id, text=messageTutor, parse_mode="html")
        logger.info(f"[NOTIFY-TUTOR] response: {response.response_id} to user: {response.tutor_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-TUTOR] response: {response.response_id} to user: {response.tutor_id} failed!. Error: {ex}")

async def proceed_selected(selected_order: SelectedDto) -> None:
    messageStudent = (
        f"<b>Вы успешно выбрали репетитора на заказ \"{selected_order.order_name}\"!</b>\n\n"
        "👀 <i>Хороших занятий!</i>"
    )

    messageTutor = (
        f"<b>Вас выбрали в качестве репетитора на заказ \"{selected_order.order_name}\"!</b>\n\n"
        "✅ <i>Продуктивных занятий!</i>"
    )

    try:
        await bot_student.send_message(chat_id=selected_order.student_telegram_id, text=messageStudent, parse_mode="html")
        logger.info(f"[NOTIFY-STUDENT] selected: {selected_order.response_id} to user: {selected_order.student_telegram_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-STUDENT] selected: {selected_order.response_id} to user: {selected_order.student_telegram_id} failed!. Error: {ex}")

    try:
        await bot_tutor.send_message(chat_id=selected_order.tutor_telegram_id, text=messageTutor, parse_mode="html")
        logger.info(f"[NOTIFY-TUTOR] selected: {selected_order.response_id} to user: {selected_order.tutor_telegram_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-TUTOR] selected: {selected_order.response_id} to user: {selected_order.tutor_telegram_id} failed!. Error: {ex}")


async def proceed_suggest(suggest_order: SuggestDto) -> None:
    message = (
        f"<b>Ученик предлагает вам заказ: {suggest_order.order_name}</b>\n\n"
        f"👀 <b>Описание:</b> {suggest_order.description}\n\n"
        f"<b>Бюджет:</b> {suggest_order.min_price} - {suggest_order.max_price}\n\n"
        "⚡ <i>Вы можете просмотреть заказ, нажав на кнопку ниже.</i>"
    )

    tutor_id = suggest_order.tutor_telegram_id

    try:
        await bot_tutor.send_message(chat_id=tutor_id, text=message, parse_mode="html", reply_markup=suggest_keyboard(suggest_order.order_id))
        logger.info(f"[NOTIFY-TUTOR] suggest order: {suggest_order.order_id} to user: {suggest_order.tutor_telegram_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-TUTOR] suggest order: {suggest_order.order_id} to user: {suggest_order.tutor_telegram_id} failed!. Error: {ex}!")

async def proceed_review(new_review: ReviewDto) -> None:
    message = (
        f"<b>Ученик оставил отзыв по заказу: {new_review.order_name}</b>\n\n"
        f"<b>Если вы занимались с учеником - подтвердите на странице профиля!</b>\n\n"
    )

    try:
        await bot_tutor.send_message(chat_id=new_review.tutor_telegram_id, text=message, parse_mode="html")
        logger.info(f"[NOTIFY-TUTOR] approved review: {new_review.order_id} to user: {new_review.tutor_telegram_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-TUTOR] approved review: {new_review.order_id} to user: {new_review.tutor_telegram_id} failed!. Error: {ex}!")

async def proceed_need_review(order: OrderDto) -> None:
    message = (
        f"<b>Вы занимались с репетитором по заказу: {order.order_name}?</b>\n\n"
        "📩 <i>Оставьте отзыв о работе - специалист будет рад узнать ваши впечатления!</i>"
    )

    try:
        await bot_student.send_message(chat_id=str(order.student_id), text=message, parse_mode="html")
        logger.info(f"[NOTIFY-STUDENT] order {order.order_id} need_review to user {order.student_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-STUDENT] order {order.order_id} need_review to user {order.student_id} failed!. Error: {ex}")

async def proceed_add_response(add_response: AddResponseDto) -> None:
    message = (
        f"<b>🎉Вам были добавлены отклики.</b>\n\n"
        f"<i>Текущее количество: {add_response.response_count}</i>"
    )

    try:
        await bot_tutor.send_message(chat_id=str(add_response.tutor_telegram_id), text=message, parse_mode="html")
        logger.info(f"[NOTIFY-TUTOR] add response message to user {add_response.tutor_telegram_id} send!")
    except Exception as ex:
        logger.error(f"[NOTIFY-TUTOR] add response message to user {add_response.tutor_telegram_id} failed!. Error: {ex}!")
