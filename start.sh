app=datahub_repository


#export goservice_port=8088
#export MYSQL_PORT_3306_TCP_ADDR=54.223.244.55
#export MYSQL_PORT_3306_TCP_PORT=3386
#export MYSQL_ENV_MYSQL_DATABASE=datahub
#export MYSQL_ENV_MYSQL_USER=root
#export MYSQL_ENV_MYSQL_PASSWORD=mysqladmin
#
#pkill -f server
export goservice_port=8088
export MYSQL_PORT_3306_TCP_ADDR=127.0.0.1
export MYSQL_PORT_3306_TCP_PORT=3306
export MYSQL_ENV_MYSQL_DATABASE=datahub
export MYSQL_ENV_MYSQL_USER=root
export MYSQL_ENV_MYSQL_PASSWORD=root


export MONGO_PORT_27017_TCP_ADDR=localhost
export MONGO_PORT_27017_TCP_PORT=27017
export MONGO_ENV_MYSQL_DATABASE=
export MONGO_ENV_MYSQL_USER=
export MONGO_ENV_MYSQL_PASSWORD=

./$app