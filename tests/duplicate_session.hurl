# second PUT must return 409

# First PUT
PUT http://localhost:8080/{{uuid}}

HTTP/1.1 201


# Duplicate PUT
PUT http://localhost:8080/{{uuid}}

HTTP/1.1 409
