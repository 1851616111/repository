curl http://127.0.0.1:8088/permission/chai -H user:chaizs@asiainfo.com
curl http://127.0.0.1:8088/permission/chai -H user:chaizs@asiainfo.com -X PUT -d '{"username":"panxy4@asiainfo.com","opt_permission":1 }'
curl http://127.0.0.1:8088/permission/chai?username=panxy3@asiainfo.com -H user:chaizs@asiainfo.com -X DELETE

curl http://127.0.0.1:8088/permission/chai/zong -H user:chaizs@asiainfo.com
curl http://127.0.0.1:8088/permission/chai/zong -H user:chaizs@asiainfo.com -X PUT -d '{"username":"panxy01@asiainfo.com"}'
curl http://127.0.0.1:8088/permission/chai/zong?username=panxy01@asiainfo.com -H user:chaizs@asiainfo.com -X DELETE