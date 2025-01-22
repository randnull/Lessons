import os
from typing import Callable

import aio_pika

class RabbitMQConsumer:
    def __init__(self, proceed_func: Callable):
        self.user: str = os.getenv('RABBITMQ_USER', "guest") # TODO что-то с кредами приудмать
        self.password: str = os.getenv('RABBITMQ_PASSWORD', 'guest')
        self.host: str = os.getenv('RABBITMQ_HOST', 'rabbitmq')
        self.port: int = int(os.getenv('RABBITMQ_PORT', 5672))
        self.queue_name: str = os.getenv('RABBITMQ_QUEUE_NAME', 'my_queue')
        self.creds: str = f"amqp://{self.user}:{self.password}@{self.host}:{self.port}/"
        self.queue: aio_pika.Queue | None = None

        self.proceed_function: Callable = proceed_func

        self.connection = None
        self.channel: aio_pika.Channel | None = None

    async def connect(self):
        self.connection = await aio_pika.connect_robust(self.creds)
        self.channel = await self.connection.channel()

        await self.channel.set_qos(prefetch_count=1)

        self.queue = await self.channel.declare_queue(self.queue_name, durable=True)

    async def disconnect(self):
        pass

    async def consume(self):
        # TODO loger
        if self.connection is None:
            raise RuntimeError("RabbitMQ connection not established")

        if self.channel is None:
            raise RuntimeError("RabbitMQ channel not established")

        if self.queue is None:
            raise RuntimeError("RabbitMQ queue not established")

        await self.queue.consume(self.proceed_function)
