from fastapi import APIRouter, Request

from aiogram.types import Update

from AnswerEngine.src.TelegramBot.bot import dp, bot


webhook_router = APIRouter()


@webhook_router.post("/webhook")
async def webhook(request: Request):
    update = Update.model_validate(await request.json(), context={"bot": bot})
    await dp.feed_update(bot, update)
