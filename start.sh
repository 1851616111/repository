app=datahub_repository


#pkill -f server
export goservice_port=8089
export MYSQL_PORT_3306_TCP_ADDR=127.0.0.1
export MYSQL_PORT_3306_TCP_PORT=3306
export MYSQL_ENV_MYSQL_DATABASE=datahub
export MYSQL_ENV_MYSQL_USER=root
export MYSQL_ENV_MYSQL_PASSWORD=root

export DB_MONGO_URL=localhost
export DB_MONGO_PORT=27017
export MONGO_ENV_MYSQL_USER=
export MONGO_ENV_MYSQL_PASSWORD=

./$app