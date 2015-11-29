curl http://10.1.235.98:8888/repositories/app -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ=="
curl http://54.223.58.0:8888/repositories/app -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ==" -x proxy.asiainfo.com:8080

curl http://10.1.235.98:8888/permission/chai -H user:chaizs@asiainfo.com

permission
-------------------------------------------------------------------------------------------------------
curl http://10.1.235.98:8888/permission/chai -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"
curl http://10.1.235.98:8888/permission/chai -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X PUT -d '{"username":"panxy4@asiainfo.com","opt_permission":1 }'
curl http://10.1.235.98:8888/permission/chai?username=panxy4@asiainfo.com -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X DELETE

curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"
curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X PUT -d '{"username":"panxy01@asiainfo.com"}'
curl http://10.1.235.98:8888/permission/chai/zong?username=panxy01@asiainfo.com -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X DELETE
-------------------------------------------------------------------------------------------------------


permission
-------------------------------------------------------------------------------------------------------
curl http://54.223.58.0:8888/permission/test007 -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -x proxy.asiainfo.com:8080
curl http://54.223.58.0:8888/permission/test007/test0007 -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -x proxy.asiainfo.com:8080
curl -X PUT  http://54.223.58.0:8888/permission/test007 -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -d '{"username":"chai@asiainfo.com","opt_permission":1 }' -x proxy.asiainfo.com:8080
curl -X PUT  http://54.223.58.0:8888/permission/test007/test0007 -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -d '{"username":"chai@asiainfo.com","opt_permission":1 }'   -x proxy.asiainfo.com:8080

curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"
curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X PUT -d '{"username":"panxy01@asiainfo.com"}'
curl http://10.1.235.98:8888/permission/chai/zong?username=panxy01@asiainfo.com -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X DELETE


curl http://127.0.0.1:8088/permission/test006 -X PUT -d '{"username":"chai@asiainfo.com"}' -H user:panxy3@asiainfo.com
curl http://127.0.0.1:8088/permission/test007 -H user:panxy3@asiainfo.com
curl http://127.0.0.1:8088/permission/test007/test007 -X PUT -d '{"username":"chai@asiainfo.com"}' -H user:panxy3@asiainfo.com
curl http://127.0.0.1:8088/permission/test007/test007 -H user:panxy3@asiainfo.com

-------------------------------------------------------------------------------------------------------

subpermission
-------------------------------------------------------------------------------------------------------

curl http://127.0.0.1:8088/repositories/test007/test007/subpermission -H user:chai@asiainfo.com
curl http://54.223.58.0:8888/repositories/test007/test0007/subpermission -H user:chai@asiainfo.com -x proxy.asiainfo.com:8080

-------------------------------------------------------------------------------------------------------


repositories
-------------------------------------------------------------------------------------------------------
curl -X PUT http://dogfood.hub.dataos.io/repositories/chai/zong  -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -d '{"comment":"123"}'
curl -X GET http://10.1.235.98:8888/repositories/chai/zong  -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"


curl -X PUT http://54.223.58.0:8888/repositories/chai/zong  -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -d '{"comment":"123"}'
curl -X GET http://54.223.58.0:8888/repositories/chai/zong

curl -X POST http://127.0.0.1:8088/repositories/test006 -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": ""}' -H user:panxy3@asiainfo.com
curl -X POST http://127.0.0.1:8088/repositories/test006/test006 -d '{"repaccesstype": "public","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "batch"},"opt": {},"owner": {},"other": {}}}' -H user:panxy3@asiainfo.com

curl -X POST http://127.0.0.1:8088/repositories/test007 -d '{"repaccesstype": "private","comment": "中国移动北京终端详情","label": ""}' -H user:panxy3@asiainfo.com
curl -X POST http://127.0.0.1:8088/repositories/test007/test007 -d '{"repaccesstype": "private","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "batch"},"opt": {},"owner": {},"other": {}}}' -H user:panxy3@asiainfo.com

curl -X POST http://54.223.58.0:8888/repositories/test007 -d '{"repaccesstype": "private","comment": "中国移动北京终端详情","label": ""}' -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -x proxy.asiainfo.com:8080
curl -X POST http://54.223.58.0:8888/repositories/test007/test0007 -d '{"repaccesstype": "private","comment": "中国移动北京终端详情","label": {"sys": {"supply_style": "batch"},"opt": {},"owner": {},"other": {}}}' -H "Authorization: Token 2b0964abd566324935ed25be65be2ad2" -x proxy.asiainfo.com:8080





-------------------------------------------------------------------------------------------------------


curl http://54.223.58.0:8888/repositories/app -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ=="
curl http://54.223.58.0:8888/permission/chai -H user:chaizs@asiainfo.com


statis
-------------------------------------------------------------------------------------------------------
curl http://54.223.58.0:8888/repositories/statis -H "Authorization:Token 81cada2f839839c12b01bdf3261bef05" -x proxy.asiainfo.com:8080
curl http://hub.dataos.io/repositories/statis -H "Authorization:Token 81cada2f839839c12b01bdf3261bef05" -x proxy.asiainfo.com:8080

label
-------------------------------------------------------------------------------------------------------
curl http://127.0.0.1:8088/repositories/app/label -d "owner.abc=100&other.name=panxy3" -X PUT -H user:panxy3@asiainfo.com
curl http://127.0.0.1:8088/repositories/liu/xu/label -d "other.name=1" -X PUT -H user:panxy3@asiainfo.com

curl -X PUT  http://127.0.0.1:8088/repositories/app/label -d "other.name=panxy3" -H user:panxy3@asiainfo.com
curl -X DELETE  http://127.0.0.1:8088/repositories/liu/xu/label?other.name=1 -H user:panxy3@asiainfo.com

curl http://10.1.235.98:8089/repositories/app0001/label -d "other.name=panxy3" -X PUT -H user:panxy3@asiainfo.com
curl http://10.1.235.98:8888/repositories/liu/xu/label -d "other.name=1" -X PUT -H "Authorization:Token ff2592ae306c68eafd89edd4b1a8e022"

curl -X DELETE  http://10.1.235.98:8088/repositories/app0001/label?other.name -H user:panxy3@asiainfo.com
curl -X DELETE  http://127.0.0.1:8088/repositories/app0001/label?other.name -H user:panxy3@asiainfo.com
curl -X DELETE  http://10.1.235.98:8088/repositories/liu/xu/label?other.name=1 -H user:panxy3@asiainfo.com


58e83462759adf3881079bd9b4dc68c81


-----------------------------------------------------------------------------------------------------
