package api

import (
	"conversational_ai/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type Server struct {
	config util.Config
	router *gin.Engine
	io     *socketio.Server
	client *openai.Client
}

func NewServer(config util.Config) (*Server, error) {
	client := openai.NewClient(config.OpenAPIKey)

	server := &Server{
		config: config,
		client: client,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping().Result()

	if err != nil {
		log.Error().Err(err).Msgf("cannot connect to redis %s", err.Error())
	}

	// Init Socket IO
	server.setupSocketIO()

	// Init Gin router
	server.setUpRouter()

	return server, nil
}

func (s *Server) setUpRouter() {
	router := gin.New()

	// Incase of a backend crush it will return a 500 error
	router.Use(util.GinRecovery())

	// Make health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Set up socket.io endpoint
	ioRoutes := router.Group("/io").Use(allowOrigin("*"))
	ioRoutes.GET("/", gin.WrapH(s.io))
	ioRoutes.POST("/", gin.WrapH(s.io))

	s.router = router
}

func (s *Server) setupSocketIO() {
	timeout := time.Duration(s.config.SocketIOPingTimeout) * time.Second
	interval := time.Duration(s.config.SocketIOPingInterval) * time.Second

	options := &engineio.Options{
		PingTimeout:  timeout,
		PingInterval: interval,
	}

	sock := socketio.NewServer(options)

	redisOpts := &socketio.RedisAdapterOptions{
		Addr:   s.config.RedisAddr,
		Prefix: s.config.RedisPrefix,
	}

	ok, err := sock.Adapter(redisOpts)
	log.Info().Msgf("Redis adapter created %v", ok)
	if condition := ok && err == nil; !condition {
		log.Error().Err(err).Msgf("cannot connect to redis %s", err.Error())
	}

	// Handle socket.io events
	sock.OnConnect("/", s.onConnect)
	sock.OnDisconnect("/", s.onDisconnect)

	/// Set up upload audio event
	sock.OnEvent("/", "stream-audio", s.streamAudio)
	sock.OnEvent("/", "audio-details", s.audioDetails)

	// Handle socket.io errors
	sock.OnError("/", s.onError)

	s.io = sock
}

func (s *Server) StartServer() error {
	port := fmt.Sprintf(":%v", s.config.HttpServerAddress)
	s.router.SetTrustedProxies([]string{"127.0.0.1"})

	// Run Socket.IO
	go func() {
		if err := s.io.Serve(); err != nil {
			log.Error().Err(err).Msg("cannot start socket.io server")
		}
	}()

	defer s.io.Close()
	return s.router.Run(port)
}
