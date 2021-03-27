# Account Management using USSD and SMS Africastalking AT

## Prerequisites

1. Go 
2. Ngrok
3. Docker 

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
