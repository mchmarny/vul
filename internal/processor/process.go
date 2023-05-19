package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/pubsub"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (h *Handler) processHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "process")

	b, err := pubsub.GetData(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "error extracting pubsub data")) //nolint:errcheck
		return
	}

	// capture image from pubsub queue
	imageURI := string(b)
	log.Debug().Str("imageURI", imageURI).Msg("processing image")

	if err := h.processImage(c.Request.Context(), imageURI); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "error processing image")) //nolint:errcheck
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": h.Version,
	})
}
