from fastapi import Depends

from AnswerEngine.common.generic_repository.repo_connection import get_repository
from AnswerEngine.src.models.dto_table.order_status_dto import OrderEngineDto


# from AnswerEngine.models.dao_table.dao import Client
# from AnswerEngine.models.dto_table.dto import ClientModel


class AnswersService:
    def __init__(self, answer_repository=Depends(get_repository(OrderEngineDto))):
        self.__answer_repository = answer_repository

    async def get_responses_by_order_id(self, request):
        return ["6b2b7f8e-c1d1-4acb-91f5-a7cc2132db72, 6b2b7f8e-c1d1-4acb-91f5-a7cc2132db73", "6b2b7f8e-c1d1-4acb-91f5-a7cc2132db72"]




        # is_client_exist = await self.__client_repository.check_by_columns(request)
        #
        # id_client = 0
        #
        # if is_client_exist is None:
        #     new_client = ClientModel.to_dto_from_request(request)
        #     id_client = await self.__client_repository.add(new_client)
        # else:
        #     id_client = is_client_exist.id_client
        #
        # return id_client
