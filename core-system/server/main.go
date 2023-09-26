package main

import (
	"bytes"
	"compress/gzip"
	"log"
	"math"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/googollee/go-socket.io"
	"net/http"
)

const (
	MaxUnackedMessages = 1000
	MessageExpiration  = time.Minute * 5
	ChannelBufferSize  = 100
	MaxRetries         = 3
)

type Message struct {
	Event string
	Data  string
	Ack   bool
}

type Subscriber struct {
	sendChan chan Message
	mu       sync.Mutex
	rdb      *redis.Client
	server   *go_socket_io.Server
}

func (s *Subscriber) handle() {
	// Assume context for Redis.
	// ctx := context.Background()

	// Handle incoming and outgoing messages.
	for msg := range s.sendChan {
		ack := sendOverSocket(msg)
		if ack {
			msg.Ack = true
		} else {
			for i := 0; i < MaxRetries; i++ {
				backoffTime := time.Duration(math.Pow(2, float64(i))) * time.Second
				time.Sleep(backoffTime)
				ack = sendOverSocket(msg)
				if ack {
					msg.Ack = true
					break
				}
			}
			if !msg.Ack {
				// Compress the message data before persisting to Redis
				var buffer bytes.Buffer
				gz := gzip.NewWriter(&buffer)
				if _, err := gz.Write([]byte(msg.Data)); err != nil {
					log.Printf("Failed to compress message data: %v\n", err)
				}
				if err := gz.Close(); err != nil {
					log.Printf("Failed to close gzip writer: %v\n", err)
				}
				compressedData := buffer.Bytes()

				// Storing compressed message in Redis is assumed here.

				log.Printf("Message not acknowledged after max retries: Event - %s, Data - %v\n", msg.Event, msg.Data)
			}
		}
	}
}

func CreateServer(s *Subscriber) {
	http.Handle("/socket.io/", s.server)
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (s *Subscriber) initializeHandlers() {
	s.server.OnConnect("/", func(so go_socket_io.Conn) error {
		so.SetContext("")
		log.Println("connected:", so.ID())
		return nil
	})

	s.server.OnEvent("/", "transcript_received", func(so go_socket_io.Conn, transcript string) string {
		log.Println("transcript received:", transcript)
		response := "Your processed data here..."
		return response
	})

	s.server.OnDisconnect("/", func(so go_socket_io.Conn, reason string) {
		log.Println("closed", reason)
	})
}

func initSubscriber() *Subscriber {
	server := go_socket_io.NewServer()
	s := &Subscriber{
		sendChan: make(chan Message, ChannelBufferSize),
		server:   server,
		// Redis client initialization is assumed here.
	}
	s.initializeHandlers()
	return s
}

func main() {
	subscriber := initSubscriber()
	CreateServer(subscriber)
	// Other initializations and function calls as necessary.
}
