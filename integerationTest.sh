#/bin/bash
Host=http://10.1.235.98:8888
#Host=http://54.223.58.0:8888
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

echo "<----------------------------- 		环境准备,清理数据 	 	    ------------------------------------->"
echo "<----------------------------- 		Start 			    	    ------------------------------------->"

echo "<-----------------------------     	删除可能存在 Tag:$Tag	    ------------------------------------->"

result=`curl -X DELETE ${Host}/repositories/$Rep/$Item/$Tag -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<-----------------------------     	删除可能存在 Item:$Item	    ------------------------------------->"
result=`curl -X DELETE ${Host}/repositories/$Rep/$Item -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<-----------------------------    	删除可能存在 Rep:$Rep	    ------------------------------------->"
result=`curl -X DELETE ${Host}/repositories/$Rep -H "Authorization:$Token" -x proxy.asiainfo.com:8080 -s`

echo "<----------------------------- 		End 			        ------------------------------------->"


