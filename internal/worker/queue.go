package worker

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (h *Handler) queueHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "queue")

	list, err := data.ListImages(c.Request.Context(), h.Pool)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error listing images")) //nolint:errcheck
		return
	}

	log.Debug().Int("count", len(list)).Msg("publishing images")

	for _, image := range list {
		if err := h.Publisher.PublishStr(c.Request.Context(), h.Config.PubSub.Topic, image); err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error publishing image to: %s", h.Config.PubSub.Topic)) //nolint:errcheck
			return
		}
	}

	log.Debug().Int("count", len(list)).Msg("published images")

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": h.Version,
	})
}
