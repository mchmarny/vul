package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "images")

	list, err := data.ListImages(c.Request.Context(), h.Pool)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error listing images")) //nolint:errcheck
		return
	}

	resp := Response[[]string]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Data:    list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func (h *Handler) imageViewHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "images_view")

	img := c.Query("img")
	log.Debug().Str("image", img).Msg("image view")

	if img == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image")) //nolint:errcheck
		return
	}

	sum, err := data.GetSummary(c.Request.Context(), h.Pool, img)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error getting image summary for %s", img)) //nolint:errcheck
		return
	}

	list, err := data.ListImageVersions(c.Request.Context(), h.Pool, img, h.Config.App.ImageVersionLimit)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image versions for %s", img)) //nolint:errcheck
		return
	}

	d := gin.H{
		"name":    h.Name,
		"version": h.Version,
		"img":     img,
		"data":    sum,
		"list":    list,
	}

	c.HTML(http.StatusOK, "image", d)
}
