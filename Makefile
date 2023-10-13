redis:
	docker run --name cai_redis --network cai-network -p 6379:6379 -d redis:7-alpine

build:
	docker build -t cai .

run:
	docker run --name cai --network cai-network -p 4400:4400 -d cai

server:
	go run main.go

# Necessary for setting up custom STUN/TURN server for webRTC
corturn:
	docker run -d -p 3478:3478 -p 3478:3478/udp -p 5349:5349 -p 5349:5349/udp -p 49152-65535:49152-65535/udp coturn/coturn
	

.PHONY: redis server build run