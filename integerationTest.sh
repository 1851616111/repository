#/bin/bash
Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
user=panxy3@asiainfo.com
passwd=88888888
admin=panxy3@asiainfo.com
admin_passwd=88888888

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
    if [ ${#basic} -e 0 ];then
        echo "no basic avaliable"
    fi

    echo $basic
}

function getToken() {
    tokenURL=${Host}/permission/mob
    token=`curl  $tokenURL -H "Authorization: Basic $Basic" -x proxy.asiainfo.com:8080 -s`
    Token=`echo $token | cut -d \" -f 4`
    if [ ${#Token} -ne 32 ];then
        echo "no token avaliable"
    fi
    Token="Token $Token"
}

function getAdminToken() {
    tokenURL=${Host}/permission/mob
        admintoken=`curl  $tokenURL -H "Authorization: Basic ZGF0YWh1YkBhc2lhaW5mby5jb206NDZjNWZjODQ5MWI5NjMyNDAxYTIwN2M3YWIwNGViMGE=" -x proxy.asiainfo.com:8080 -s`
    AdminToken=`echo $admintoken | cut -d \" -f 4`
    if [ ${#AdminToken} -ne 32 ];then
        echo "no admintoken avaliable"
    fi
    AdminToken="Token $AdminToken"
}

function chkResult() {
    msg=`echo $1 | cut -d "," -f 2 | cut -d ":" -f 2`
    if [ "${msg:1:2}" != "OK" ];then
        echo "$1"
    fi
}

$(getBasic) $user $passwd
$(getBasic) $admin $admin_passwd
echo "----------> $Basic"
getToken
echo "Token : $Token"

getAdminToken
echo "AdminToken : $AdminToken"


Rep=Repository_$RANDOM
Item=Dataitem_$RANDOM
Tag=Tag_$RANDOM
Label=Label_$RANDOM



echo "<----------------------------- 		环境准备,清理数据 	 	    ------------------------------------->"
echo "<----------------------------- 		Start 			    	    ------------------------------------->"

echo "<-----------------------------     	删除可能存在 Tag:$Tag	    ------------------------------------->"

result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<-----------------------------     	删除可能存在 Item:$Item	    ------------------------------------->"
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<-----------------------------    	删除可能存在 Rep:$Rep	    ------------------------------------->"
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<----------------------------- 		End 			        ------------------------------------->"


