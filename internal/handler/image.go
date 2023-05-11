package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
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
