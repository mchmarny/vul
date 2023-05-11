package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
)

func (h *Handler) imageVersionExposureHandler(c *gin.Context) {
	var criteria vul.ListImageVersionExposureRequest
	if err := c.ShouldBindJSON(&criteria); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "binding image version exposure request"))
		return
	}

	if criteria.Image == "" || criteria.Digest == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image or its digest"))
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_version_exposure", map[string]string{
		"image": criteria.Image,
	})

	list, err := data.ListImageVersionExposures(c.Request.Context(), h.Pool, criteria.Image, criteria.Digest)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image version exposures for %s", criteria.Image))
		return
	}

	resp := Response[map[string][]*vul.ListDigestExposureItem]{
		Version:  h.Version,
		Created:  time.Now().UTC(),
		Criteria: criteria,
		Data:     list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
