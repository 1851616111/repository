app=datahub_repository


export goservice_port=8088
export MYSQL_PORT_3306_TCP_ADDR=54.223.244.55
export MYSQL_PORT_3306_TCP_PORT=3386
export MYSQL_ENV_MYSQL_DATABASE=datahub
export MYSQL_ENV_MYSQL_USER=root
export MYSQL_ENV_MYSQL_PASSWORD=mysqladmin
#
##pkill -f server
#export goservice_port=8088
#export MYSQL_PORT_3306_TCP_ADDR=127.0.0.1
#export MYSQL_PORT_3306_TCP_PORT=3306
#export MYSQL_ENV_MYSQL_DATABASE=datahub
#export MYSQL_ENV_MYSQL_USER=root
#export MYSQL_ENV_MYSQL_PASSWORD=root

./$app