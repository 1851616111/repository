app=datahub_repository

#pkill -f server
export API_SERVER=54.223.58.0
export API_PORT=8888

export goservice_port=8088

export MONGO_PORT_27017_TCP_ADDR=54.223.58.0
export MONGO_PORT_27017_TCP_PORT=27019
export MONGO_ENV_MYSQL_USER=
export MONGO_ENV_MYSQL_PASSWORD=
export MQ_KAFKA_ADDR=10.1.235.98
export MQ_KAFKA_PORT=9092

if [ $1 -z ];then
    ./$app
else
    $1 $2 $3 $4 $5 $6
fi

