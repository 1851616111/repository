#/bin/bash
Host=http://10.1.235.98:8888
#Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
user=panxy3@asiainfo.com
passwd=88888888
admin=datahub@asiainfo.com
admin_passwd=BDXLDPdatahub

Basic=""
AdminBasic=""
Token=""
AdminToken=""

function getBasic() {

    user=$1
    password=$2
    
    if [ -z "$user" ];then
        echo username null
	exit
    fi
    if [ -z "$password" ];then
        echo password null
	exit
    fi
    
    pw=`echo -n $2 | md5sum`
    tmp=${pw:0:32}
    basic=`echo -n $1:$tmp |base64`
    if [ "${basic}" = "" ];then
        echo "no basic avaliable"
        exit
    fi

    echo $basic
}

function getToken() {

   basic=$1

   if [ -z "$basic" ];then
       echo basic null
   exit
   fi

    tokenURL=${Host}/permission/mob
    token=`curl  $tokenURL -H "Authorization: Basic $basic" -x proxy.asiainfo.com:8080 -s`
    token=`echo $token | cut -d \" -f 4`

    if [ ${#token} -ne 32 ];then
        echo "no token avaliable"
        exit
    fi

    echo "Token $token"

}

function chkResult() {
    msg=`echo $1 | cut -d "," -f 2 | cut -d ":" -f 2`
    if [ "${msg:1:2}" != "OK" ];then
        echo "$1"
    fi
}

Basic=$(getBasic $user $passwd)
AdminBasic=$(getBasic $admin $admin_passwd)

Token=$(getToken $Basic)
echo "Token : $Token"

AdminToken=$(getToken $AdminBasic)
echo "AdminToken : $AdminToken"

Rep=Repository_$RANDOM
Item=Dataitem_$RANDOM
Tag=Tag_$RANDOM
Label=Label_$RANDOM
NewLabel=NewLabel_$RANDOM

echo "<----------------------------- 		环境准备,清理数据"
echo "<----------------------------- 		Start"
echo "<-----------------------------     	删除可能存在 Tag:$Tag"
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
echo "<-----------------------------     	删除可能存在 Item:$Item"
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
echo "<-----------------------------    	删除可能存在 Rep:$Rep"
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
echo "<----------------------------- 		删除可能存在 SelectLabel 	     "
result=`curl -X DELETE ${Host}/select_labels/$Label -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
echo "<----------------------------- 		End"

echo "1.-----------------------------> 		【拥有者】【新增】Rep        ($Rep)"
result=`curl -X POST ${Host}/repositories/$Rep -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "2.-----------------------------> 		【拥有者】【新增】Item       ($Rep/$Item) 	     "
result=`curl -X POST ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "api"},"opt": {},"owner": {},"other": {}},"price":[{"times": 1000,"money": 5,"expire":30},{"times": 10000,"money": 45,"expire":30},{"times":100000,"money": 400,"expire":30}]}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "3.-----------------------------> 		【拥有者】【新增】Tag        ($Rep/$Item/$Tag) 	     "
result=`curl -X POST ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "4.-----------------------------> 		【拥有者】【更新】Rep        ($Rep) 	     "
result=`curl -X PUT ${Host}/repositories/$Rep -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": ""}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "5.-----------------------------> 		【拥有者】【更新】Item	     ($Rep/$Item) 	     "
result=`curl -X PUT ${Host}/repositories/$Rep/$Item -d '{"repaccesstype": "private","comment": "update中国移动北京终端详情","label": {"sys": {"supply_style": "api"},"opt": {},"owner": {},"other": {}}}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "6.-----------------------------> 		【拥有者】【更新】Tag        ($Rep/$Item/$Tag) 	     "
result=`curl -X PUT ${Host}/repositories/$Rep/$Item/$Tag -d '{"comment":"update this is a tag"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "7.-----------------------------> 		【管理员】【新增】SelectLabel 	     "
result=`curl -X POST ${Host}/select_labels/$Label -d '{"order": 1,"icon":"path1"}' -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "8.-----------------------------> 		【管理员】【查询】SelectLabel 	     "
result=`curl -X GET ${Host}/select_labels  -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "9.-----------------------------> 		【管理员】【更新】SelectLabel 	     "
result=`curl -X PUT ${Host}/select_labels/$Label -d '{"order": 1,"icon":"path12", "newlabelname":"$NewLabel"}' -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "10.-----------------------------> 		【管理员】将Item添加至精选 	     "
result=`curl -X POST ${Host}/selects/$Rep/$Item -d "select_labels=$Label&order=1" -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "11.-----------------------------> 		【管理员】【更新】Item的精选 	     "
result=`curl -X PUT ${Host}/selects/$Rep/$Item -d "select_labels=$Label&order=88" -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "12.-----------------------------> 		【任意】【查询】精选的Item 	     "
result=`curl  ${Host}/selects?select_labels=$Label  -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "13.-----------------------------> 		【拥有者】【删除】Item 	     "
result=`curl -X DELETE ${Host}/selects/$Rep/$Item  -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "14.-----------------------------> 		【拥有者】【删除】SelectLabel 	     "
result=`curl -X DELETE ${Host}/select_labels/$Label -H "Authorization:$AdminToken" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "15.-----------------------------> 		【任意】【查询】Rep 	     "
 result=`curl -X GET ${Host}/repositories/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "16.-----------------------------> 		【任意】【查询】Item 	     "
result=`curl -X GET ${Host}/repositories/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "17.-----------------------------> 		【任意】【查询】Tag 	     "
result=`curl -X GET ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "18.-----------------------------> 		【拥有者】【新增】某用户添加至Rep白名单 	     "
result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "19.-----------------------------> 		【拥有者】【新增】某用户添加至Item白名单 	     "
result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com"}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "20.-----------------------------> 		【拥有者】【更新】某用户在白名单中的权限 	     "
 result=`curl -X PUT ${Host}/permission/$Rep -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "21.-----------------------------> 		【拥有者】【更新】某用户在白名单中的权限 	     "
result=`curl -X PUT ${Host}/permission/$Rep/$Item  -d '{"username":"chai@asiainfo.com","opt_permission":1}' -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "22.-----------------------------> 		【拥有者】【查询】Rep白名单列表 	     "
result=`curl -X GET ${Host}/permission/$Rep  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "23.-----------------------------> 		【拥有者】【查询】Item白名单列表 	     "
result=`curl -X GET ${Host}/permission/$Rep/$Item  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "24.-----------------------------> 		【拥有者】【删除】Item白名单用户	      "
result=`curl -X DELETE ${Host}/permission/$Rep/$Item?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "25.-----------------------------> 		【拥有者】【删除】Rep白名单用户	     "
result=`curl -X DELETE ${Host}/permission/$Rep?username=chai@asiainfo.com  -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "26.-----------------------------> 		【拥有者】【删除】Tag 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "27.-----------------------------> 		【拥有者】【删除】Item 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "28.-----------------------------> 		【拥有者】【删除】Rep 	     "
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "29.-----------------------------> 		【任意】【查询】Search 关键字 (mobile)	     "
result=`curl -X GET ${Host}/search?text=mobile -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

echo "30.-----------------------------> 		【任意】【查询】Search 关键字 (mobile空格pc)	     "
result=`curl -X GET ${Host}/search?text=mobile -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`
chkResult $result

