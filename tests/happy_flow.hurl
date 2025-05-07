# Happy path: PUT → POST → GET (SSE) → DELETE

# Create an interactive shell session
PUT http://localhost:8080/{{uuid}}

HTTP/1.1 201

# Run a command inside the shell
POST http://localhost:8080/{{uuid}}
Content-Type: text/plain

```
echo hello && echo done && exit
```

HTTP/1.1 202

# Stream output via SSE
GET http://localhost:8080/{{uuid}}

HTTP/1.1 200
[Asserts]
header "Content-Type" contains "text/event-stream"
body contains "data:hello"
body contains "data:done"

# Terminate the session
DELETE http://localhost:8080/{{uuid}}

HTTP/1.1 204
