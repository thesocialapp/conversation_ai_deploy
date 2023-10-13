package api

import (
	"bytes"

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

func (server *Server) streamAudio(io socketio.Conn, data map[string]interface{}) {
	/// We use this to append the file chunks to a buffer
	var fileBuffer bytes.Buffer

	// Process the file chunk data
	log.Info().Msgf("Client file chunk: %s", data)
	chunkData, ok := data["data"].([]byte)
	if !ok {
		io.Emit("fileChunk", "fileparse error")
		return
	}

	fileName, ok := data["fileName"].(string)
	if !ok {
		io.Emit("fileChunk", "fileName error")
		return
	}

	_, err := fileBuffer.Write(chunkData)
	if err != nil {
		// Tell the user the file is too large
		if err == bytes.ErrTooLarge {
			io.Emit("fileChunk", "file too large")
			return
		}
		io.Emit("fileChunk", "fileBuffer error")
		return
	}
	io.Emit("transcriptionResult", "ok "+fileName)
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
