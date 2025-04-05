from fastapi import Depends

from AnswerEngine.common.generic_repository.repo_connection import get_repository
# from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto
#
from AnswerEngine.src.models.dao_table.dao import TutorTagDao

# from AnswerEngine.models.dao_table.dao import Client
# from AnswerEngine.models.dto_table.dto import ClientModel


class TutorTagService:
    def __init__(self, tags_repository=Depends(get_repository(TutorTagDao))):
        self.__tags_repository = tags_repository

    async def set_tutor_tags(self, request):
        pass

    # async def get_responses_by_order_id(self, request):
    #     return ["6b2b7f8e-c1d1-4acb-91f5-a7cc2132db72, 6b2b7f8e-c1d1-4acb-91f5-a7cc2132db73", "6b2b7f8e-c1d1-4acb-91f5-a7cc2132db72"]
    #
    #
