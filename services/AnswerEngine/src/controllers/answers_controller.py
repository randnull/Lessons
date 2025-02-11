from fastapi import APIRouter, Depends
from fastapi.responses import JSONResponse

from AnswerEngine.src.services.answers_service import AnswersService

from typing import Annotated

answers_router = APIRouter(prefix="/responses", tags=["responses"])

AnswersServiceBase = Annotated[AnswersService, Depends(AnswersService)]

@answers_router.post("/{order_id}", tags=["responses"])
async def get_responses(order_id: int, answers_service: AnswersServiceBase):
    """
    Get responses by order_id
    """
    # try:
    responses = await answers_service.get_responses_by_order_id(order_id)
    # except AgreementNotFound:
    #     return JSONResponse(content={"message": "Agreement not found"}, status_code=404)

    # return JSONResponse(content={"message": "Agreement was closed"}, status_code=200)
    return JSONResponse(responses)