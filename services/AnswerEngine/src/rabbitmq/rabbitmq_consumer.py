from AnswerEngine.common.rabbitmq.consumer import RabbitMQConsumer
from AnswerEngine.src.rabbitmq.rabbitmq_functions import user_func

OrderConsumer = RabbitMQConsumer(user_func)


# Task exception was never retrieved

# ne/src/rabbitmq/bot_functions.py", line 10, in proceed_order
#     await bot.send_message(chat_id=order_create.chat_id, title=order_create.order_title)
#           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
# TypeError: Bot.send_message() got an unexpected keyword argument 'title'
# Task exception was never retrieved
# future: <Task finished name='Task-19' coro=<consumer() done, defined at /Users/kirillgorunov/PycharmProjects/PythonProject/.venv/lib/python3.12/site-packages/aio_pika/queue.py:28> exception=TypeError("Bot.send_message() got an unexpected keyword argument 'title'")>
# Traceback (most recent call last):
#   File "/Users/kirillgorunov/PycharmProjects/PythonProject/.venv/lib/python3.12/site-packages/aio_pika/queue.py", line 34, in consumer
#     return await create_task(callback, message)
#            ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
#   File "/Users/kirillgorunov/PycharmProjects/PythonProject/AnswerEngine/src/rabbitmq/rabbitmq_functions.py", line 17, in user_func