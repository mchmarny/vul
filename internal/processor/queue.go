package processor

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) queueHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "queue")

	// TODO:
	// - get images from DB
	// - push each image to PubSub topic

	resp := Response[[]string]{
		Version: h.Version,
		Created: time.Now().UTC(),
	}

	c.IndentedJSON(http.StatusOK, resp)
}
