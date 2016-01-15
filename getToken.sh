Host=http://10.1.235.98:8888
#Host=http://54.223.58.0:8888
#Host=http://hub.dataos.io/api
#user=panxy3@asiainfo.com
#passwd=88888888
user=shendf@asiainfo.com
passwd=99999999
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
   echo "Basic : $basic"
   if [ -z "$basic" ];then
       echo basic null
   exit
   fi

    tokenURL=${Host}
    #itoken=`curl  $tokenURL -H "Authorization: Basic $basic" -x proxy.asiainfo.com:8080 -s`
    token=`curl  $tokenURL -H "Authorization: Basic $basic" -s`
    token=`echo $token | cut -d \" -f 4`

    if [ ${#token} -ne 32 ];then
        echo "no token avaliable"
    fi

    echo "Token $token"

}

Basic=$(getBasic $user $passwd)
AdminBasic=$(getBasic $admin $admin_passwd)

Token=$(getToken $Basic)
echo "Token : $Token"

AdminToken=$(getToken $AdminBasic)
echo "AdminToken : $AdminToken"