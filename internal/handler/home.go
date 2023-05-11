package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
)

func (h *Handler) homeHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "home")

	list, err := data.ListImages(c.Request.Context(), h.Pool)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error listing images"))
		return
	}

	sum, err := data.GetSummary(c.Request.Context(), h.Pool, "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error getting image summary"))
		return
	}

	d := gin.H{
		"name":    h.Name,
		"version": h.Version,
		"images":  list,
		"summary": sum,
	}

	c.HTML(http.StatusOK, "home", d)
}
