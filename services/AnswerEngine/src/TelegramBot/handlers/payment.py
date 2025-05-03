from aiogram import Bot
from aiogram.types import Message, LabeledPrice, PreCheckoutQuery, CallbackQuery

from AnswerEngine.src.TelegramBot.keyboards.keyboards import payment_keyboard
from AnswerEngine.src.config.settings import settings
from AnswerEngine.src.functions.payment_functions import add_tutor_responses



async def send_invoice_handler(message: Message):
    await message.answer(
        text="Выберите количество откликов:",
        reply_markup=payment_keyboard()
    )


async def process_subscription_callback(callback_query: CallbackQuery, bot: Bot):
    subscription_type = callback_query.data

    if subscription_type == "sub_5":
        amount = 5
        description = "Оплатить 10⭐"
    elif subscription_type == "sub_10":
        amount = 10
        description = "Оплатить 20⭐"
    elif subscription_type == "sub_15":
        amount = 15
        description = "Оплатить 30⭐"
    elif subscription_type == "sub_30":
        amount = 30
        description = "Оплатить 50⭐"
    else:
        await callback_query.answer("Недостустимое количетсво!")
        return

    prices = [LabeledPrice(label="XTR", amount=amount)] # amount
    await bot.send_invoice(
        chat_id=callback_query.from_user.id,
        title="Покупка откликов",
        description=description,
        prices=prices,
        provider_token="",
        payload=f"subscription:{subscription_type}",
        currency="XTR",
    )
    await callback_query.answer()


async def pre_checkout_handler(pre_checkout_query: PreCheckoutQuery):
    payload = pre_checkout_query.invoice_payload
    total_amount = pre_checkout_query.total_amount
    currency = pre_checkout_query.currency
    tutor_id = pre_checkout_query.from_user.id

    if not payload.startswith("subscription:"):
        await pre_checkout_query.answer(
            ok=False,
            error_message="Неверные данные платежа."
        )
        return

    sub_type = int(payload.split("_")[1])
    #
    if sub_type not in [5, 10, 15, 30]:
        await pre_checkout_query.answer(
            ok=False,
            error_message="Ошибка валидации."
        )
        return

    if currency != "XTR":
        await pre_checkout_query.answer(
            ok=False,
            error_message="Поддерживается только XTR."
        )
        return

    responses, status = await add_tutor_responses(tutor_id, total_amount)

    if not status:
        await pre_checkout_query.answer(ok=False, error_message="Пожалуйста, попробуйте позже")

    await pre_checkout_query.answer(ok=True)


async def success_payment_handler(message: Message):
    payload = message.successful_payment.invoice_payload
    response_count = payload.split("_")[1]

    await message.answer(text=f"🥳 Вы успешно приобрели отклики: {response_count} уже на вашем балансе!")


async def pay_support_handler(message: Message):
    await message.answer(text=f"Если у вас возникли какие-либо проблемы - свяжитесь с поддержкой: {settings.SUPPORT_CHANNEL}")

