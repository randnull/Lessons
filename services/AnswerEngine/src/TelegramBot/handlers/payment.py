from aiogram import Bot
from aiogram.types import Message, LabeledPrice, PreCheckoutQuery, CallbackQuery

from AnswerEngine.src.TelegramBot.keyboards.keyboards import payment_keyboard


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
        amount = 30
        description = "Оплатить 30⭐"
    elif subscription_type == "sub_30":
        amount = 30
        description = "Оплатить 50⭐"
    else:
        await callback_query.answer("Недостустимое количетсво!")
        return

    prices = [LabeledPrice(label="XTR", amount=amount)]
    await bot.send_invoice(
        chat_id=callback_query.from_user.id,
        title="Подписка на канал",
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

    if not payload.startswith("subscription:"):
        await pre_checkout_query.answer(
            ok=False,
            error_message="Неверные данные платежа."
        )
        return

    sub_type = payload.split(":")[1]

    if sub_type == "sub_5" and total_amount != int(payload.split("_")[1]):
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

    # if requests.get("https://lessonsmy.tech").status_code != 205:
    #     await pre_checkout_query.answer(
    #         ok=False,
    #         error_message="Сайт недоступен, попробуйте позже."
    #     )
    #     return

    await pre_checkout_query.answer(ok=True)


async def success_payment_handler(message: Message):
    payload = message.successful_payment.invoice_payload
    response_count = payload.split("_")[1]

    await message.answer(text=f"🥳 Вы успешно приобрели отклики: {response_count} уже на вашем балансе!")
