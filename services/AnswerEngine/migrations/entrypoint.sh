#!/bin/sh
cat <<EOF > /tmp/answer.properties
classpath: /app
url: jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}
username: ${DB_USER}
password: ${DB_PASSWORD}
changeLogFile: changelog.xml
EOF

liquibase --defaultsFile=/tmp/answer.properties update