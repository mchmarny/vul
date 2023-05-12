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

func (h *Handler) imageTimelineHandler(c *gin.Context) {
	img := c.Query("img")
	log.Debug().Str("image", img).Msg("image version")

	if img == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image")) //nolint:errcheck
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_timeline", map[string]string{
		"image": img,
	})

	since := time.Now().UTC().
		AddDate(0, 0, -h.Config.App.ImageTimelineDays).
		Format(time.DateOnly)

	list, err := data.ListImageTimelines(c.Request.Context(), h.Pool, img, since)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image timeline for %s", img)) //nolint:errcheck
		return
	}

	resp := Response[[]*vul.ImageTimeline]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Criteria: map[string]interface{}{
			"image": img,
			"since": since,
		},
		Data: list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
