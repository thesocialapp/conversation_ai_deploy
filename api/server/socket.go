package api

import (
	"context"
	"encoding/base64"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

func (server *Server) onConnect(conn socketio.Conn) error {
	log.Info().Msgf("Client connected: %s", conn.ID())

	return nil
}

func (server *Server) onDisconnect(io socketio.Conn, reason string) {
	log.Error().Msgf("Client disconnected: %s", reason)
}

func (server *Server) audioDetails(io socketio.Conn, data string) {
	log.Info().Msgf("Client audio details: %s", data)
	io.Emit("audioResponse", "ok ")
}

type AudioData struct {
	FileName string `msgpack:"name"`
	Audio    string `msgpack:"audio"`
}

func (server *Server) streamAudio(io socketio.Conn, data string) {
	log.Info().Msg("Client audio stream")
	// Convert data string to buffer base64
	parsedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error().Msgf("Error decoding base64: %s", err.Error())
		io.Emit("transcriptionResult", "error decoding"+err.Error())
		return
	}

	var audioData AudioData
	// // Parse the message using message pack
	if err := msgpack.Unmarshal(parsedBytes, &audioData); err != nil {
		log.Error().Msgf("Error parsing message pack: %s", err.Error())
		io.Emit("transcriptionResult", "error buffer "+err.Error())
		return
	}

	// Save the audio file
	fileName := audioData.FileName + ".ogg"
	log.Info().Msgf("Saving audio file: %s", fileName)

	r, err := server.rClient.Publish(context.Background(), "audio", audioData.Audio).Result()
	if err != nil {
		log.Error().Msgf("Error publishing message: %s", err.Error())
		return
	}
	log.Info().Msgf("Client audio stream: %v", r)
	/// Try to save the audio file

	io.Emit("transcriptionResult", "ok "+audioData.FileName)
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
