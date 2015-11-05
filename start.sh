app=goservice

#pkill -f server
export goservice_port=8080
export MYSQL_PORT_3306_TCP_ADDR=127.0.0.1
export MYSQL_PORT_3306_TCP_PORT=3306
export MYSQL_ENV_MYSQL_DATABASE=datahub
export MYSQL_ENV_MYSQL_USER=root
export MYSQL_ENV_MYSQL_PASSWORD=root

./$app