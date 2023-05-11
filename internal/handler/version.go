package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
)

func (h *Handler) imageVersionHandler(c *gin.Context) {
	var criteria vul.ImageRequest
	if err := c.ShouldBindJSON(&criteria); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "error binding image version request"))
		return
	}

	if criteria.Image == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty image"))
		return
	}

	h.Meter.RecordOneWithLabels(c.Request.Context(), "image_version", map[string]string{
		"image": criteria.Image,
	})

	list, err := data.ListImageVersions(c.Request.Context(), h.Pool, criteria.Image)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrapf(err, "error listing image versions for %s", criteria.Image))
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
