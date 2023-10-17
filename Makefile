redis:
	docker run --name cai_redis --network xai-network -p 6379:6379 -d redis:7-alpine

build:
	docker build -t cai -f Dockerfile_go .

run:
	docker run --name cai --network xai-network -p 4400:4400 -d cai

build_py:
	docker build -t cai_py -f Dockerfile_python .

run_py:
	docker run --name cai_py --network xai-network -p 4401:4401 -d cai_py

server:
	go run main.go

# Necessary for setting up custom STUN/TURN server for webRTC
corturn:
	docker run -d -p 3478:3478 -p 3478:3478/udp -p 5349:5349 -p 5349:5349/udp -p 49152-65535:49152-65535/udp coturn/coturn
	

.PHONY: redis server build build_py run run_py corturn