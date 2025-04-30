from AnswerEngine.common.rabbitmq.consumer import RabbitMQConsumer
from AnswerEngine.src.rabbitmq.rabbitmq_functions import new_order_func, response_func, suggest_func, \
    tutors_change_tags_func, selected_order_func, review_func

OrderConsumer = RabbitMQConsumer(new_order_func, "new_orders")
ResponseConsumer = RabbitMQConsumer(response_func, "order_response")
SuggestConsumer = RabbitMQConsumer(suggest_func, "suggest_order")
TagsChangeConsumer = RabbitMQConsumer(tutors_change_tags_func, "tutors_tags_change")
SelectedConsumer = RabbitMQConsumer(selected_order_func, "selected_orders")
ReviewConsumer = RabbitMQConsumer(review_func, "new_review")
