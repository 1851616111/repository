app=datahub_repository

#pkill -f server
export goservice_port=8080
export MYSQL_PORT_3306_TCP_ADDR=10.1.235.96
export MYSQL_PORT_3306_TCP_PORT=3306
export MYSQL_ENV_MYSQL_DATABASE=datahub
export MYSQL_ENV_MYSQL_USER=datahub
export MYSQL_ENV_MYSQL_PASSWORD=datahub

./$app