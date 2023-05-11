package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageVersionExposureHandler(c *gin.Context) {
	var criteria query.ListImageVersionExposureRequest
	if err := c.ShouldBindJSON(&criteria); err != nil {
		log.Error().Err(err).Msg("error binding request")
		c.JSON(http.StatusBadRequest, ErrInvalidRequest)
		c.Abort()
		return
	}

	if criteria.Image == "" || criteria.Digest == "" {
		log.Error().Msg("empty image or its digest")
		c.JSON(http.StatusBadRequest, ErrInvalidRequest)
		c.Abort()
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_version_exposure", map[string]string{
		"image": criteria.Image,
	})

	list, err := data.ListImageVersionExposures(c.Request.Context(), h.Pool, criteria.Image, criteria.Digest)
	if err != nil {
		log.Error().Err(err).Msgf("error listing image version exposures for %s", criteria.Image)
		c.JSON(http.StatusInternalServerError, ErrInternal)
		c.Abort()
		return
	}

	resp := Response[map[string][]*query.ListDigestExposureItem]{
		Version:  h.Version,
		Created:  time.Now().UTC(),
		Criteria: criteria,
		Data:     list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
