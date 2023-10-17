package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) UploadAudioFile(ctx *gin.Context) {
	/// Respond with a json
	ctx.JSON(http.StatusOK, "We are working")
}

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
