#/bin/bash
Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
Token=""
function getToken() {
    tokenURL=${Host}/permission/mob
    token=`curl  $tokenURL -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ==" -x proxy.asiainfo.com:8080`
    Token=`echo $token | cut -d \" -f 4`
    if [ ${#Token} -ne 32 ];then
        echo "no token avaliable"
    fi
    Token="Token $Token"
}

function chkResult() {
    msg=`echo $1 | cut -d "," -f 2 | cut -d ":" -f 2`
    if [ "${msg:1:2}" != "OK" ];then
        echo "$1 xx"
    fi
}

getToken
echo "Token : $Token"

Rep=random$RANDOM
Item=random$RANDOM
Tag=random$RANDOM


result=`curl -X POST ${Host}/repositories/$Rep -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X POST ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "batch"},"opt": {},"owner": {},"other": {}}}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X POST ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result


result=`curl -X PUT ${Host}/repositories/$Rep -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": {"sys": {"supply_style": "api"},"opt": {},"owner": {},"other": {}}}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"update this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X GET ${Host}/repositories/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X GET ${Host}/repositories/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X GET ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X GET ${Host}/permission/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X GET ${Host}/permission/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X DELETE ${Host}/permission/$Rep/$Item?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X DELETE ${Host}/permission/$Rep?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

result=`curl -X POST ${Host}/select_labels/testlabel -d '{"order": 1,"icon":"path1"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
echo $result

#result=`curl -X PUT ${Host}/select_labels/testlabel -d '{"order": 1,"icon":"path12", "newlabelname":"2015股市"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
#chkResult $result

result=`curl -X GET ${Host}/select_labels  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

#result=`curl -X Delete ${Host}/select_labels/2015股市 -d '{"order": 1,"icon":"path1"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
#chkResult $result