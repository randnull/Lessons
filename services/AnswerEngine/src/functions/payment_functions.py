import grpc

from AnswerEngine.src.config.settings import settings
from AnswerEngine.src.gRPC import tutors_pb2, tutors_pb2_grpc

grpc_host = settings.GRPC_HOST
grpc_port = settings.GRPC_PORT

async def add_tutor_responses(tutor_id: int, response_count: int = 1) -> tuple[int, bool]:
    async with grpc.aio.insecure_channel(f'{grpc_host}:{grpc_port}') as grpc_channel:
        grpc_stub = tutors_pb2_grpc.UserServiceStub(grpc_channel)

        try:
            response = await grpc_stub.AddResponsesToTutor(
                tutors_pb2.AddResponseToTutorRequest(
                    tutor_id=tutor_id,
                    response_count=response_count
                )
            )
            return response.response_count, response.success
        except grpc.RpcError as e:
            print(f"gRPC error while adding responses for tutor {tutor_id}: {e}")
            return 0, False