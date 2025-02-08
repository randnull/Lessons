# from fastapi import Depends
#
# from AnswerEngine.common.generic_repository.repo_connection import get_repository
# from AnswerEngine.models.dao_table.dao import Client
# from AnswerEngine.models.dto_table.dto import ClientModel
#
#
# class AnswersService:
#     def __init__(self, client_repository=Depends(get_repository(Client))):
#         self.__client_repository = client_repository
#
#     async def get_client_id(self, request):
#         is_client_exist = await self.__client_repository.check_by_columns(request)
#
#         id_client = 0
#
#         if is_client_exist is None:
#             new_client = ClientModel.to_dto_from_request(request)
#             id_client = await self.__client_repository.add(new_client)
#         else:
#             id_client = is_client_exist.id_client
#
#         return id_client
