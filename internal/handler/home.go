package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
)

func (h *Handler) homeViewHandler(c *gin.Context) {
	h.Recorder.One(c.Request.Context(), "home")

	list, err := data.ListImages(c.Request.Context(), h.Pool)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error listing images")) //nolint:errcheck
		return
	}

	sum, err := data.GetSummary(c.Request.Context(), h.Pool, "", "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "error getting image summary")) //nolint:errcheck
		return
	}

	d := gin.H{
		"name":    h.Name,
		"version": h.Version,
		"images":  list,
		"data":    sum,
	}

	c.HTML(http.StatusOK, "home", d)
}
