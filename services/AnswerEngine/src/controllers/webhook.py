from fastapi import APIRouter, Request

from aiogram.types import Update

from AnswerEngine.src.TelegramBot.botStudent import dp_student, bot_student
from AnswerEngine.src.TelegramBot.botTutor import dp_tutor, bot_tutor

webhook_router = APIRouter()


@webhook_router.post("/webhook/student")
async def student_webhook(request: Request):
    update = Update.model_validate(await request.json(), context={"bot": bot_student})
    await dp_student.feed_update(bot_student, update)
    return {"status": "ok"}

@webhook_router.post("/webhook/tutor")
async def tutor_webhook(request: Request):
    update = Update.model_validate(await request.json(), context={"bot": bot_tutor})
    await dp_tutor.feed_update(bot_tutor, update)
    return {"status": "ok"}