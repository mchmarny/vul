package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageSummaryHandler(c *gin.Context) {
	img := c.Query("img")
	log.Debug().Str("image", img).Msg("image summary")

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_summary", map[string]string{
		"image": img,
	})

	list, err := data.GetSummary(c.Request.Context(), h.Pool, img)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error getting image summary"))
		return
	}

	resp := Response[*vul.SummaryItem]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Data:    list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
