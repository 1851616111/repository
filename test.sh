curl http://10.1.235.98:8888/repositories/app -H "Authorization: Basic Y2hhaXpzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ===="
curl http://54.223.58.0:8888/repositories/app -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ=="
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
curl http://54.223.58.0:8888/permission/chai -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"
curl http://54.223.58.0:8888/permission/chai -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X PUT -d '{"username":"panxy4@asiainfo.com","opt_permission":1 }'
curl http://54.223.58.0:8888/permission/chai?username=panxy4@asiainfo.com -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X DELETE

curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"
curl http://10.1.235.98:8888/permission/chai/zong -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X PUT -d '{"username":"panxy01@asiainfo.com"}'
curl http://10.1.235.98:8888/permission/chai/zong?username=panxy01@asiainfo.com -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -X DELETE
-------------------------------------------------------------------------------------------------------



repositories
-------------------------------------------------------------------------------------------------------
curl -X PUT http://dogfood.hub.dataos.io/repositories/chai/zong  -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -d '{"comment":"123"}'
curl -X GET http://10.1.235.98:8888/repositories/chai/zong  -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd"


curl -X PUT http://54.223.58.0:8888/repositories/chai/zong  -H "Authorization: Token 5c69d07b6143bc4c443564203b3704fd" -d '{"comment":"123"}'
curl -X GET http://54.223.58.0:8888/repositories/chai/zong
-------------------------------------------------------------------------------------------------------


curl http://54.223.58.0:8888/repositories/app -H "Authorization: Basic cGFueHkzQGFzaWFpbmZvLmNvbTo4ZGRjZmYzYTgwZjQxODljYTFjOWQ0ZDkwMmMzYzkwOQ=="
curl http://54.223.58.0:8888/permission/chai -H user:chaizs@asiainfo.com


statis
-------------------------------------------------------------------------------------------------------
curl http://54.223.58.0:8888/repositories/statis -H "Authorization:Token 81cada2f839839c12b01bdf3261bef05"
