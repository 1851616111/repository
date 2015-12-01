Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
Token=""
AdminToken=""
function getToken() {
    tokenURL=${Host}/permission/mob
    token=`curl  $tokenURL -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ==" -x proxy.asiainfo.com:8080`
    Token=`echo $token | cut -d \" -f 4`
    if [ ${#Token} -ne 32 ];then
        echo "no token avaliable"
    fi
    Token="Token $Token"
}

function getAdminToken() {
    tokenURL=${Host}/permission/mob
        admintoken=`curl  $tokenURL -H "Authorization: Basic ZGF0YWh1YkBhc2lhaW5mby5jb206NDZjNWZjODQ5MWI5NjMyNDAxYTIwN2M3YWIwNGViMGE=" -x proxy.asiainfo.com:8080`
    AdminToken=`echo $admintoken | cut -d \" -f 4`
    if [ ${#AdminToken} -ne 32 ];then
        echo "no admintoken avaliable"
    fi
    AdminToken="Token $AdminToken"
}

function chkResult() {
    msg=`echo $1 | cut -d "," -f 2 | cut -d ":" -f 2`
    if [ "${msg:1:2}" != "OK" ];then
        echo "$1 xx"
    fi
}

getToken
echo "Token : $Token"

getAdminToken
echo "AdminToken : $AdminToken"


Rep=random$RANDOM
Item=random$RANDOM
Tag=random$RANDOM
Label=random$RANDOM

echo 创建repository
result=`curl -X POST ${Host}/repositories/$Rep -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 创建dataitem
result=`curl -X POST ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "batch"},"opt": {},"owner": {},"other": {}}}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 创建tag
result=`curl -X POST ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 更新repository
result=`curl -X PUT ${Host}/repositories/$Rep -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 更新dataitem
result=`curl -X PUT ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": {"sys": {"supply_style": "api"},"opt": {},"owner": {},"other": {}}}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 更新tag
result=`curl -X PUT ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"update this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 添加精选栏目
result=`curl -X POST ${Host}/select_labels/$Label -d '{"order": 1,"icon":"path1"}' -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

#echo 精选栏目重命名
#result=`curl -X PUT ${Host}/select_labels/$Label -d '{"order": 1,"icon":"path12", "newlabelname":"2016"}' -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
#chkResult $result

echo 查询精选栏目
result=`curl -X GET ${Host}/select_labels  -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

echo 将dataitem添加至精选
result=`curl -X POST ${Host}/selects/$Rep/$Item -d "select_labels=$Label&order=100" -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

echo 修改dataitem的精选
result=`curl -X PUT ${Host}/selects/$Rep/$Item -d "select_labels=$Label&order=88" -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询精选的dataitem
result=`curl  ${Host}/selects?select_labels=$Label  -x proxy.asiainfo.com:8080`
chkResult $result

echo 删除dataitem
result=`curl -X DELETE ${Host}/selects/$Rep/$Item  -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

echo 删除精选栏目
result=`curl -X DELETE ${Host}/select_labels/$Label -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询repsitory
result=`curl -X GET ${Host}/repositories/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询dataitem
result=`curl -X GET ${Host}/repositories/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询tag
result=`curl -X GET ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 将某用户添加至repository的白名单
result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 将某用户添加至dataitem的白名单
result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 更新某用户在白名单中的权限
result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result


echo 更新某用户在白名单中的权限
result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询自己repository白名单列表
result=`curl -X GET ${Host}/permission/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 查询自己dataitem白名单列表
result=`curl -X GET ${Host}/permission/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 将某用户从item白名单列表中删除
result=`curl -X DELETE ${Host}/permission/$Rep/$Item?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 将某用户从rep白名单列表中删除
result=`curl -X DELETE ${Host}/permission/$Rep?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 删除tag
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 删除item
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

echo 删除rep
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080`
chkResult $result

