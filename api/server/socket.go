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

func (server *Server) audioDetails(io socketio.Conn, data map[string]interface{}) {
	log.Info().Msgf("Client audio details: %s", data["file"])
	io.Emit("audioResponse", "ok ")
}

func (server *Server) rtcOffer(io socketio.Conn, data map[string]interface{}) {
	sdp, err := server.peerConn.ProcessOffer(data["offer"].(string))
	if err != nil {
		log.Error().Msgf("Client rtcOffer error: %s", err.Error())
		io.Emit("rtcResponse", "error")
		return
	}
	io.Emit("rtcResponse", sdp)
	io.Emit("rtcResponse", "ok ")
}

func (server *Server) streamAudio(io socketio.Conn, data map[string]interface{}) {
	/// We use this to append the file chunks to a buffer
	var fileBuffer bytes.Buffer

	// Process the file chunk data
	log.Info().Msgf("Client file chunk: %s", data["chunk"])
	chunkData, ok := data["chunk"].([]byte)
	log.Info().Msgf("Client file chunk: %s", len(chunkData))

	if !ok {
		io.Emit("audioResponse", "fileparse error")
		return
	}

	fileName, ok := data["blob-name"].(string)
	if !ok {
		io.Emit("audioResponse", "fileName error")
		return
	}

	_, err := fileBuffer.Write(chunkData)
	if err != nil {
		// Tell the user the file is too large
		if err == bytes.ErrTooLarge {
			io.Emit("audioResponse", "file too large")
			return
		}
		io.Emit("audioResponse", "fileBuffer error")
		return
	}
	io.Emit("transcriptionResult", "ok "+fileName)
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
