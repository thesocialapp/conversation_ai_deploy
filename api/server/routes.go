package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) UploadAudioFile(ctx *gin.Context) {
	/// Respond with a json
	ctx.JSON(http.StatusOK, "We are working")
}
