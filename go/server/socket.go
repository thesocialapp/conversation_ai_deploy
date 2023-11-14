package server

import (
	"context"
	"encoding/base64"
	"os"

	socketio "github.com/googollee/go-socket.io"
	el "github.com/haguro/elevenlabs-go"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
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

	io.Emit("transcriptionResult", "**waiting**")

	log.Info().Msgf("Transcription: %s", resp.Text)

	completion, err := server.llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{
			Content: "Hello! I'm here to help with questions related to airplanes. If you have any inquiries about aircraft, aviation, or related topics, feel free to ask, and I'll do my best to provide you with accurate information." +
				"However, if your question is not related to airplanes, " +
				" I'll still try to assist you as politely as possible." +
				" Please keep your questions respectful and on-topic. How can I assist you today? " +
				"This message sets the expectation that the AI is focused on airplanes but also emphasizes " +
				"that it will politely respond to other questions while gently encouraging users to stick to the topic",
		},
		schema.HumanChatMessage{
			Content: resp.Text,
		},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		// log.Info().Msgf("Here it is %s", string(chunk))
		//  Need an io.Write to write the audio to and then read it later on
		//  to send it to the client
		io.Emit("transcriptionResult", string(chunk))
		return nil
	}))
	if err != nil {
		log.Error().Msgf("Error publishing message: %s", err.Error())
	}

	// r, err := server.rClient.Publish(context.Background(), "audio", resp.Text).Result()
	// if err != nil {
	// 	log.Error().Msgf("Error publishing message: %s", err.Error())
	// }
	audio, err := server.elClient.TextToSpeech("pNInz6obpgDQGcFmaJgB", el.TextToSpeechRequest{
		Text:    completion.Content,
		ModelID: "eleven_monolingual_v1",
	})
	if err != nil {
		log.Error().Msgf("Error converting text to speech: %s", err.Error())
	}
	io.Emit("audio_response", base64.StdEncoding.EncodeToString(audio))
	log.Info().Msgf("Client audio stream: %v", completion.GetContent())
	/// Try to save the audio file
}

func (server *Server) onError(io socketio.Conn, err error) {
	log.Error().Msgf("Client error: %s", err.Error())
}
