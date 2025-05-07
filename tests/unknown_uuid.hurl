# unknown ID should always give 404

GET http://localhost:8080/{{unknown}}

HTTP/1.1 404


POST http://localhost:8080/{{unknown}}
Content-Type: text/plain

```
noop
```

HTTP/1.1 404


DELETE http://localhost:8080/{{unknown}}

HTTP/1.1 404
