## GO CRUD APP
A Simple CRUD APP.

Using: 
* net/http package for routing
* mysql for database
* sqlc for sql code generation

## Run locally

### Prerequisites
- docker
- go (>1.22)

### Steps
* ``docker compose up --build -d``

### Endpoints:
* POST http://localhost:8080/user
* GET http://localhost:8080/user?limit=10&offset=0
* GET  http://localhost:8080/user/1
* DELETE http://localhost:8080/user/1
