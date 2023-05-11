package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
)

func (h *Handler) imageSummaryHandler(c *gin.Context) {
	var img string

	if c.Request.Method == http.MethodPost {
		var criteria vul.ImageRequest
		if err := c.BindJSON(&criteria); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "error binding image summary request"))
			return
		}
		img = criteria.Image
	}

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
