package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageVersionHandler(c *gin.Context) {
	var criteria vul.ListImageVersionRequest
	if err := c.ShouldBindJSON(&criteria); err != nil {
		log.Error().Err(err).Msg("error binding request")
		c.JSON(http.StatusBadRequest, ErrInvalidRequest)
		c.Abort()
		return
	}

	if criteria.Image == "" {
		log.Error().Msg("empty image")
		c.JSON(http.StatusBadRequest, ErrInvalidRequest)
		c.Abort()
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_version", map[string]string{
		"image": criteria.Image,
	})

	list, err := data.ListImageVersions(c.Request.Context(), h.Pool, criteria.Image)
	if err != nil {
		log.Error().Err(err).Msgf("error listing image versions for %s", criteria.Image)
		c.JSON(http.StatusInternalServerError, ErrInternal)
		c.Abort()
		return
	}

	resp := Response[map[string][]*vul.ListImageSourceItem]{
		Version:  h.Version,
		Created:  time.Now().UTC(),
		Criteria: criteria,
		Data:     list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
