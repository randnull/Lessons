from fastapi import APIRouter, Depends
from fastapi.responses import JSONResponse

from AnswerEngine.src.services.tags_service import TutorTagService

from typing import Annotated

tags_router = APIRouter(prefix="/tutor/tags", tags=["tutor"])

TagsServiceBase = Annotated[TutorTagService, Depends(TutorTagService)]

@tags_router.get("/{tutor_id}", tags=["tutor"])
async def get_responses(tutor_id: str, tags_service: TagsServiceBase):
    """
    Set new tutor tags
    """
    # try:
    # responses = await tags_service.get_responses_by_order_id(order_id)
    # except AgreementNotFound:
    #     return JSONResponse(content={"message": "Agreement not found"}, status_code=404)

    # return JSONResponse(content={"message": "Agreement was closed"}, status_code=200)
    # return JSONResponse(responses)