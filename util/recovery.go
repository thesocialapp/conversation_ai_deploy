package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GinRecovery is a gin middleware that recovers from any panics and writes a 500 if there was one.
func GinRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal().Msgf("panic recovered :%v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
		}()
		// Continue processing the request
		ctx.Next()
	}
}
