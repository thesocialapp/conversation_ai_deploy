redis:
	docker run --name redis --network cai-network -p 6379:6379 -d redis:7-alpine

build:
	docker build -t cai .

run:
	docker run --name cai --network cai-network -p 4400:4400 -d cai

server:
	go run backend/main.go
	

.PHONY: redis server build run