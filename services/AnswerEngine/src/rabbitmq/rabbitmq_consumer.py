from AnswerEngine.common.rabbitmq.consumer import RabbitMQConsumer
from AnswerEngine.src.rabbitmq.rabbitmq_functions import new_order_func, response_func, suggest_func

OrderConsumer = RabbitMQConsumer(new_order_func, "new_orders")
ResponseConsumer = RabbitMQConsumer(response_func, "order_response")
SuggestConsumer = RabbitMQConsumer(suggest_func, "suggest_order")
