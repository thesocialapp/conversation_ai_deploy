package api

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

func (server *Server) onConnect(conn socketio.Conn) error {

	log.Info().Msgf("Client connected: %s", conn.ID())
	return nil
}

func (server *Server) onDisconnect(io socketio.Conn, reason string) {
	log.Error().Msgf("Client disconnected: %s", reason)
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
