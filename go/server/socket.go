package server

import (
	"context"
	"encoding/base64"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/vmihailenco/msgpack/v5"
)

func (server *Server) onConnect(conn socketio.Conn) error {
	log.Info().Msgf("Client connected: %s", conn.ID())
	ctx := context.Background()

	/// After a working connection we listen for audio responses
	/// from eleven labs
	subChan := server.rClient.Subscribe(ctx, "audio_response").Channel()
	/// Run a goroutine to listen for messages
	go func() {
		for msg := range subChan {
			// Convert the payload from base64 to bytes
			// and send it to the client
			audioByte := base64.StdEncoding.EncodeToString([]byte(msg.Payload))
			conn.Emit("audio_response", audioByte)
		}
	}()

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
	// Convert data string to buffer base64 under utf-8
	parsedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error().Msgf("Error decoding base64: %s", err.Error())
		io.Emit("transcriptionResult", "error decoding"+err.Error())

	}

	var audioData AudioData
	// // Parse the message using message pack
	if err := msgpack.Unmarshal(parsedBytes, &audioData); err != nil {
		log.Error().Msgf("Error parsing message pack: %s", err.Error())
		io.Emit("transcriptionResult", "error buffer "+err.Error())
	}

	// Save the audio file
	fileName := audioData.FileName + ".ogg"
	log.Info().Msgf("Saving audio file: %s", fileName)

	// Convert base64 to bytes
	audioBytes, err := base64.StdEncoding.DecodeString(audioData.Audio)
	if err != nil {
		log.Error().Msgf("Error decoding base64: %s", err.Error())
		io.Emit("transcriptionResult", "error decoding: "+err.Error())
	}

	/// Create the temp file from the audio bytes and io using os
	f, err := os.CreateTemp("", "audio-*.ogg")
	if err != nil {
		log.Error().Msgf("Error creating temp file: %s", err.Error())
		io.Emit("transcriptionResult", "error creating temp file"+err.Error())
		return
	}

	defer f.Close()
	defer os.Remove(f.Name())

	// Write the audio bytes to the temp file
	if _, err := f.Write(audioBytes); err != nil {
		log.Error().Msgf("Error writing to temp file: %s", err.Error())
		io.Emit("transcriptionResult", "error writing to temp file"+err.Error())

	}

	if err != nil {
		log.Error().Msgf("Error creating temp file: %s", err.Error())
		io.Emit("transcriptionResult", "error creating temp file"+err.Error())

	}

	// reader := bytes.NewReader(audioBytes)
	// Audio transcription req
	req := openai.AudioRequest{
		FilePath: f.Name(),
		Model:    openai.Whisper1,
	}
	ctx := context.Background()

	// Send the audio to openai
	// and get the transcription
	resp, err := server.client.CreateTranscription(ctx, req)
	if err != nil {
		log.Err(err).Msgf("Error creating transcription: %s", err.Error())

	}

	log.Info().Msgf("Transcription: %s", resp.Text)

	r, err := server.rClient.Publish(context.Background(), "audio", resp.Text).Result()
	if err != nil {
		log.Error().Msgf("Error publishing message: %s", err.Error())

	}
	log.Info().Msgf("Client audio stream: %v", r)
	/// Try to save the audio file

	io.Emit("transcriptionResult", "ok "+audioData.FileName)
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
