
app=datahub_repository

#pkill -f server
export API_SERVER=10.1.235.98
export API_PORT=8888

export goservice_port=8088

export CONSUL_SERVER=10.1.235.98
export CONSUL_DNS_PORT=8600
export kafka_service_name=datahub_kafka
export mongo_service_name=datahub_repository_mongo
export ADMIN_API_USERNAME=repository_9bacb3c5dt@asiainfo.com
export ADMIN_API_USER_PASSWORD=repouser888

if [ $1 -z ];then
    ./$app
else
    $1 $2 $3 $4 $5 $6
fi

