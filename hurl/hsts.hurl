GET http://localhost:8080/api/hsts?url=google.com

HTTP 200
[Asserts]
jsonpath "$.error" not exists
