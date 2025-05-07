# any call after deletion ⇒ 404

# Ensure the session is gone (idempotent)
DELETE http://localhost:8080/{{uuid}}


# GET after delete
GET http://localhost:8080/{{uuid}}

HTTP/1.1 404


# POST after delete
POST http://localhost:8080/{{uuid}}
Content-Type: text/plain

```
echo test
```

HTTP/1.1 404
