package server

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dslipak/pdf"
	"github.com/gin-gonic/gin"
)

// Create a new error const for failing to parse file
var errParsingFile = fmt.Errorf("failed to parse file")
var errOpeningFile = fmt.Errorf("failed to open file")
var errReadingFile = fmt.Errorf("failed to read file")

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

	tempFile, err := os.CreateTemp("", "file-*.pdf")
	if err != nil {
		fmt.Printf("Error creating temporary file %s", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errParsingFile))
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Write the file buffer to the temp file
	if _, err := tempFile.Write(fileBuffer); err != nil {
		fmt.Printf("Error writing to temporary file %s", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errParsingFile))
		return
	}

	// Create an inmemory file from multipart.FileHeader
	r, err := pdf.Open(tempFile.Name())
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(errOpeningFile))
		return
	}

	reader, err := r.GetPlainText()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errReadingFile))
		return
	}

	// Read the file
	b, err := io.ReadAll(reader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errReadingFile))
		return
	}
	fmt.Printf("File content: %s\n", string(b))
	// Send back the name of the file back to the client
	ctx.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
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

func getFileMIMEType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read a chunk to determine the file type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the file position
	src.Seek(0, io.SeekStart)

	// Detect the MIME type based on the file content
	return http.DetectContentType(buffer), nil
}
