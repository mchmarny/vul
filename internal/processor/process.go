package processor

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) processHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "process")

	// TODO:
	// - capture image from pubsub queue
	// scan using local commands
	// - save results to DB

	resp := Response[[]string]{
		Version: h.Version,
		Created: time.Now().UTC(),
	}

	c.IndentedJSON(http.StatusOK, resp)
}
