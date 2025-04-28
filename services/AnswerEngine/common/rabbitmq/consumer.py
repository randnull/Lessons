import os
from typing import Callable

import aio_pika

from AnswerEngine.src.logger.logger import logger

from AnswerEngine.src.config.settings import settings

class RabbitMQConsumer:
    def __init__(self, proceed_func: Callable, queue_name: str):
        self.user: str = settings.MQUSER
        self.password: str = settings.MQPASSWORD
        self.host: str = settings.MQHOST
        self.port: int = settings.MQPORT

        self.queue_name: str = queue_name

        self.creds: str = f"amqp://{self.user}:{self.password}@{self.host}:{self.port}/"
        self.queue: aio_pika.Queue | None = None

        self.proceed_function: Callable = proceed_func

        self.connection = None
        self.channel: aio_pika.Channel | None = None

        logger.info(f"RabbitMQ Consumer Initiated. Queue Name: {self.queue_name}")

    async def connect(self):
        self.connection = await aio_pika.connect_robust(self.creds)
        self.channel = await self.connection.channel()

        await self.channel.set_qos(prefetch_count=1)

        self.queue = await self.channel.declare_queue(self.queue_name, durable=True)

    async def disconnect(self):
        if self.channel and not self.channel.is_closed:
            await self.channel.close()
            logger.info("RabbitMQ channel closed.")
        if self.connection and not self.connection.is_closed:
            await self.connection.close()
            logger.info("RabbitMQ connection closed.")

    async def consume(self):
        if self.connection is None:
            logger.error(f"RabbitMQConsumer {self.connection} is None")
            raise RuntimeError("RabbitMQ connection not established")

        if self.channel is None:
            logger.error(f"RabbitMQConsumer {self.channel} is None")
            raise RuntimeError("RabbitMQ channel not established")

        if self.queue is None:
            logger.error(f"RabbitMQConsumer {self.queue} is None")
            raise RuntimeError("RabbitMQ queue not established")

        logger.debug(f"RabbitMQConsumer {self.queue} consumed")
        await self.queue.consume(self.proceed_function)
