# Account Management using USSD and SMS Africastalking AT

## Prerequisites

1. Go 
2. Ngrok
3. Docker 

## Getting started

### Running the database for the first time

1. Start the docker compose service
```sh
docker-compose up
```

2. Start the connection to database to setup the tables for the first time
```sh
psql -h <network_ip_address> -U d1r3ct0r -p 5431 -W at_reg_db -f cmd.sql
```

## How to run

1. Start the docker compose services
```sh
docker-compose up
```
**Please make sure you have the docker network specified in the docker-compose.yml**

2. Start the go code
```sh
go run *.go
```
3. Setup the ngrok endpoint and update in the *ui/templates/login.html* and *ussd.go* files 
```sh
ngrok http 8083
```
4. Enjoy and update where you see fit.
