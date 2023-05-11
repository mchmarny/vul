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

func (h *Handler) imageVersionHandler(c *gin.Context) {
	img := c.Query("img")
	log.Debug().Str("image", img).Msg("image version")

	if img == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image"))
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_version", map[string]string{
		"image": img,
	})

	list, err := data.ListImageVersions(c.Request.Context(), h.Pool, img)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image versions for %s", img))
		return
	}

	resp := Response[map[string][]*vul.ListImageSourceItem]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Criteria: map[string]interface{}{
			"image": img,
		},
		Data: list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
