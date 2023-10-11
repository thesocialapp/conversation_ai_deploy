package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {
	router := gin.New()
	server := socketio.NewServer(nil)

	redisOpts := &socketio.RedisAdapterOptions{
		Host:   "localhost:6379",
		Prefix: "io",
	}

	ok, err := server.Adapter(redisOpts)
	if condition := ok && err == nil; !condition {
		logger.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Println("connected:", s.ID())
		return nil
	})

	router.GET("/io/*any", gin.WrapH(server))
	router.POST("/io/*any", gin.WrapH(server))

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server failed to start: %v", err)
		}
	}()
	// Defer closing the server
	defer server.Close()

	if err := router.Run(":8080"); err != nil {
		logger.Fatal(err)
	}
}
