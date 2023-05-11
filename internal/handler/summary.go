package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageSummaryHandler(c *gin.Context) {
	var img string

	if c.Request.Method == http.MethodPost {
		var criteria vul.ImageRequest
		if err := c.BindJSON(&criteria); err != nil {
			log.Error().Err(err).Msg("error binding image summary request")
			c.JSON(http.StatusBadRequest, ErrInvalidRequest)
			c.Abort()
			return
		}
		img = criteria.Image
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_summary", map[string]string{
		"image": img,
	})

	list, err := data.GetSummary(c.Request.Context(), h.Pool, img)
	if err != nil {
		log.Error().Err(err).Msg("error getting image summary")
		c.JSON(http.StatusInternalServerError, ErrInternal)
		c.Abort()
		return
	}

	resp := Response[*vul.SummaryItem]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Data:    list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
