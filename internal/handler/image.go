package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "images")

	list, err := data.ListImages(c.Request.Context(), h.Pool)
	if err != nil {
		log.Error().Err(err).Msg("error listing images")
		c.JSON(http.StatusInternalServerError, ErrInternal)
		c.Abort()
		return
	}

	resp := Response[[]*query.ListImageItem]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Data:    list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
