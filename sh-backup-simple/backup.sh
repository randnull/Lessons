#!/bin/sh

ORDER_DB="postgresql"
USERS_DB="postgresql-users"
NOTIFICATION_DB="postgresql-answer"

ORDER_USER={{ ORDER_DB_NAME }}
USERS_USER={{ USERS_DB_NAME }}
NOTIFICATION_USER={{ ANSWER_DB_NAME }}

ORDER_TABLE={{ ORDER_DB_USER }}
USERS_TABLE={{ USERS_DB_USER }}
NOTIFICATION_TABLE={{ ANSWER_DB_USER }}

docker exec -t "$USERS_DB" pg_dump -U "$USERS_USER" "$USERS_TABLE" > /tmp/users_dump.sql
docker exec -t "$ORDER_DB" pg_dump -U "$ORDER_USER" "$ORDER_TABLE" > /tmp/orders_dump.sql
docker exec -t "$NOTIFICATION_DB" pg_dump -U "$NOTIFICATION_USER" "$NOTIFICATION_TABLE" > /tmp/notification_dump.sql

aws s3 cp /tmp/users_dump.sql s3://{{ backet_name_users }}--endpoint-url {{ endpoint }}
aws s3 cp /tmp/orders_dump.sql s3://{{ backet_name_orders }}--endpoint-url {{ endpoint }}
aws s3 cp /tmp/notification_dump.sql s3://{{ backet_name_notifications }}--endpoint-url {{ endpoint }}
