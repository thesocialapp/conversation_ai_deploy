package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new error const for failing to parse file
var errParsingFile = fmt.Errorf("failed to parse file")

func (s *Server) UploadAudioFile(ctx *gin.Context) {

	// Limit the file size to 8MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 8<<20)
	// Get the file from the request
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	allowedFileTypes := []string{".pdf", ".doc", ".docx", ".csv"}

	// Ensure we have a valid file type
	if !isValidFileType(file, allowedFileTypes) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid file type")))
		return
	}

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errParsingFile))
		return
	}

	// Convert the file buffer to base64
	fileBase64 := base64.StdEncoding.EncodeToString(fileBuffer)

	// Send the file to py using redis pubsub
	r, err := s.rClient.Publish(ctx, "file-document", fileBase64).Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Handle the subcribe option to ensure that the file was processed
	if r == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("failed to process file")))
		return
	}

	resultChannel := make(chan string, 1)

	// Once the py microservice process it will send a message back here to tell us that processing
	// is complete and we can now show a success message to the client
	// We set up a go routine to listen for the message from redis
	go s.subscribe("file-result", func(payload []byte) {
		resultChannel <- string(payload)
	})

	select {
	case result := <-resultChannel:
		if result == "success" {
			ctx.JSON(http.StatusOK, successResponse(result))
		} else {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("failed to process file")))
		}
	case <-time.After(10 * time.Second):
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("timeout waiting for file processing")))
	}
}

// Handle pubsub events from Redis to client
func (s *Server) ListenForPubSub(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Simulate sending events periodically
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data := "This is a sample message."
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		case <-c.Writer.CloseNotify():
			return
		}
	}
}

func getSrcFile(file *multipart.FileHeader) (multipart.File, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return src, nil
}

func getFileBuffer(file *multipart.FileHeader) ([]byte, error) {
	src, err := getSrcFile(file)
	if err != nil {
		return nil, err
	}

	// Read a chunk to determine the file type
	var buffer []byte
	buffer, err = io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func isValidFileType(file *multipart.FileHeader, allowedTypes []string) bool {
	// Get the file extension
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))

	// Filter the file extensions
	for _, allowedType := range allowedTypes {
		if fileExtension == allowedType {
			return true
		}
	}

	return false
}

// getFileMIMEType returns the MIME type of the file
// func getFileMIMEType(file *multipart.FileHeader) (string, error) {
// 	src, err := file.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer src.Close()

// 	// Read a chunk to determine the file type
// 	buffer := make([]byte, 512)
// 	_, err = src.Read(buffer)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Reset the file position
// 	src.Seek(0, io.SeekStart)

// 	// Detect the MIME type based on the file content
// 	return http.DetectContentType(buffer), nil
// }
