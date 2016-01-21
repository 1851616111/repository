
app=datahub_repository

#pkill -f server
export API_SERVER=10.1.235.98
export API_PORT=8888

export goservice_port=8088

export CONSUL_SERVER=10.1.235.98
export CONSUL_DNS_PORT=8600
export Service_Name_Kafka=datahub_kafka

if [ $1 -z ];then
    ./$app
else
    $1 $2 $3 $4 $5 $6
fi

