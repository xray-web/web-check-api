GET http://localhost:8080/api/linked-pages?url=google.com

HTTP 200
[Asserts]
jsonpath "$.error" not exists
