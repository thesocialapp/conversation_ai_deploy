redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

server:
	go run backend/main.go
	

.PHONY: redis server