# The demo project on Golang 

[![License](https://img.shields.io/badge/license-GPLv3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.html)

License: GPL v3. [GPL License](http://www.gnu.org/licenses)

[![Build Status](https://travis-ci.com/Tsovak/rest-api-demo.svg?branch=master)](https://travis-ci.com/Tsovak/rest-api-demo)
[![Coverage](https://codecov.io/gh/Tsovak/rest-api-demo/branch/master/graph/badge.svg)](https://codecov.io/gh/Tsovak/rest-api-demo)
[![GolangCI](https://golangci.com/r/github.com/Tsovak/rest-api-demo)](https://golangci.com/r/github.com/Tsovak/rest-api-demo)

## Introduction

This is the demo project which provides simple endpoints. Project store based on Postgres. 

The project contains unit tests, integration tests, docker-compose support and swagger-ui for API.

### Tests 

Unit tests can be run without any external utils. <br>
Integration tests need a running docker service or docker.sock.

### Requirements
    1. Go 1.13
    2. Running docker service or docker.sock
    3. Docker-compose 
    4. Make 
    5. Postgres 12.0 and greater 
    
## Build
   
The project has a Dockerfile. It's simple to build by Docker. <br>
Run the command in the project base dir 
`docker build -t rest-api-demo .` or `make install_deps build` 

You able to run all necessary services via docker-compose. 
`docker-compose up` <br>
Docker-compose will build the Dockerfile, setup the Postgres, run the migrations, start the application, prepare the swagger-ui 

## Deploy 

For deploying the application you need to run the migration before. 
Migration uses the same config file as the main application. <br>
`./migrate -dir ./scripts/migrations -init` or use make target `make migrate` <br>
If the command exit code is 0 it means that the migration went successfully. 

## Documentation

### Code Documentation 

Run  `godoc -http=:6060` in project directory. <br>
Read go docs on http://localhost:6060/pkg/github.com/tsovak/rest-api-demo/.

### API Documentation

API docs available by swagger. 
Just deploy using docker-compose and open `http://localhost:8888/#/`. Or open https://app.swaggerhub.com/apis-docs/Tsovak/go-rest-api/1.0.0 

## Examples

### Create an account 

Request 
```
curl --location --request POST 'http://localhost:8080/accounts' \
--header 'Content-Type: application/json' \
--data-raw '  {
    "balance": 125,
    "currency": "RU",
    "name": "alice1"
  }'
```
Response
```
{
    "id": 1,
    "name": "alice1",
    "currency": "RU",
    "balance": 125
}
```


### Get accounts  

Request 
```
curl -X GET 'http://localhost:8080/accounts'
```
Response
```
[
    {
        "ID": 1,
        "Name": "alice1",
        "Currency": "RU",
        "Balance": 125
    },
    {
        "ID": 2,
        "Name": "Alice",
        "Currency": "RU",
        "Balance": 100500
    },
    {
        "ID": 34,
        "Name": "Bob",
        "Currency": "USD",
        "Balance": 1000
    }
]
```


### Create payment

Request 
```
curl --location --request POST 'http://localhost:8080/payments' \
--header 'Content-Type: application/json' \
--data-raw '{
	"amount":1,
	"to_id":"39",
	"from_id":"40"
}'
```
Response
```
[
    {
        "id": 1,
        "amount": -1,
        "to_account": "39",
        "from_account": "40",
        "direction": "outgoing"
    },
    {
        "id": 2,
        "amount": 1,
        "to_account": "39",
        "from_account": "40",
        "direction": "incoming"
    }
]
```



### Get account payments

Request 
```
curl -X GET 'http://localhost:8080/accounts/39/payments'
```
Response
```
[
    {
        "id": 1,
        "amount": -1,
        "to_account": "39",
        "from_account": "40",
        "direction": "outgoing"
    },
    {
        "id": 2,
        "amount": 1,
        "to_account": "39",
        "from_account": "40",
        "direction": "incoming"
    },
    {
        "id": 36,
        "amount": -50,
        "to_account": "39",
        "from_account": "40",
        "direction": "outgoing"
    },
    {
        "id": 37,
        "amount": 50,
        "to_account": "39",
        "from_account": "40",
        "direction": "incoming"
    }
]
```


