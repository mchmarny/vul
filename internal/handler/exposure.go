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

func (h *Handler) imageVersionExposureHandler(c *gin.Context) {
	img := c.Query("img")
	dig := c.Query("dig")
	log.Debug().
		Str("image", img).
		Str("digest", dig).
		Msg("image version exposure")

	if img == "" || dig == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image or its digest")) //nolint:errcheck
		return
	}

	h.Recorder.OneWithLabels(c.Request.Context(), "image_version_exposure", map[string]string{
		"image":  img,
		"digest": dig,
	})

	list, err := data.ListImageVersionExposures(c.Request.Context(), h.Pool, img, dig)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image version exposures for %s@%s", img, dig)) //nolint:errcheck
		return
	}

	resp := Response[*vul.ImageDigestExposures]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Criteria: map[string]interface{}{
			"image":  img,
			"digest": dig,
		},
		Data: list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func (h *Handler) imageVersionExposureViewHandler(c *gin.Context) {
	h.Recorder.One(c.Request.Context(), "image_version_exposure_view")

	img := c.Query("img")
	dig := c.Query("dig")
	log.Debug().
		Str("image", img).
		Str("digest", dig).
		Msg("image version exposure view")

	if img == "" || dig == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image or its digest")) //nolint:errcheck
		return
	}

	sum, err := data.GetSummary(c.Request.Context(), h.Pool, img, dig)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error getting image summary for %s", img)) //nolint:errcheck
		return
	}

	list, err := data.ListImageVersionExposures(c.Request.Context(), h.Pool, img, dig)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image version exposures view for %s@%s", img, dig)) //nolint:errcheck
		return
	}

	// add fixed count
	sum.FixedCount = list.FixedCount

	d := gin.H{
		"name":    h.Name,
		"version": h.Version,
		"img":     img,
		"dig":     dig,
		"data":    sum,
		"list":    list,
	}

	c.HTML(http.StatusOK, "exposure", d)
}
