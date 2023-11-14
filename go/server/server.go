package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	util "github.com/thesocialapp/conversation-ai/go/util"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/haguro/elevenlabs-go"
	el "github.com/haguro/elevenlabs-go"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	langOpenAI "github.com/tmc/langchaingo/llms/openai"
)

type Server struct {
	config   util.Config
	router   *gin.Engine
	io       *socketio.Server
	client   *openai.Client
	rClient  *redis.Client
	llm      *langOpenAI.Chat
	elClient *el.Client
}

func NewServer(config util.Config) (*Server, error) {
	client := openai.NewClient(config.OpenAPIKey)

	// Set up langchain open ai
	llm, err := langOpenAI.NewChat(
		langOpenAI.WithToken(config.OpenAPIKey),
		langOpenAI.WithAPIVersion("v1"),
	)
	if err != nil {
		log.Error().Err(err).Msgf("cannot create langchain openai client %s", err.Error())
	}

	// Set up elevenlabs
	elClient := elevenlabs.NewClient(context.Background(), config.ElevenLabsAPIKey, 30*time.Second)

	server := &Server{
		config:   config,
		client:   client,
		llm:      llm,
		elClient: elClient,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		log.Error().Err(err).Msgf("cannot connect to redis %s", err.Error())
	}

	/// Set up redis client
	server.rClient = rdb

	// Init Socket IO
	server.setupSocketIO()

	// Init Gin router
	server.setUpRouter()

	return server, nil
}

func (s *Server) subscribeToAudioResponse() {
	ctx := context.Background()

	/// After a working connection we listen for audio responses
	/// from eleven labs
	subChan := s.rClient.Subscribe(ctx, "audio_response").Channel()
	/// Run a goroutine to listen for messages
	go func() {
		for msg := range subChan {
			// Convert the payload from base64 to bytes
			// and send it to the client
			audioByte := base64.StdEncoding.EncodeToString([]byte(msg.Payload))
			s.io.BroadcastToNamespace("/", "audio_response", audioByte)
			// conn.Emit("audio_response", audioByte)
		}
	}()

}

func (s *Server) setUpRouter() {
	/// Change gin mode to production if in production mode
	if s.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	// Rest routes
	rest := router.Group("/v1")
	rest.GET("/audio", s.UploadAudioFile)

	s.router = router
}

func (s *Server) setupSocketIO() {
	// // timeout := time.Duration(s.config.SocketIOPingTimeout) * time.Second
	// // interval := time.Duration(s.config.SocketIOPingInterval) * time.Second

	// options := &engineio.Options{
	// 	PingTimeout:  timeout,
	// 	PingInterval: interval,
	// }
	// options := &engineio.Options{
	// 	PingTimeout:  timeout,
	// 	PingInterval: interval,
	// }

	sock := socketio.NewServer(nil)

	redisOpts := &socketio.RedisAdapterOptions{
		Addr:   s.config.RedisAddr,
		Prefix: s.config.RedisPrefix,
	}

	ok, err := sock.Adapter(redisOpts)
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
