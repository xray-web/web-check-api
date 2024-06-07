GET http://localhost:8080/api/firewall?url=google.com

HTTP 200
[Asserts]
jsonpath "$.error" not exists
