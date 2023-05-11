package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) homeHandler(c *gin.Context) {
	h.Meter.RecordOne(c.Request.Context(), "home")

	d := gin.H{
		"name":    h.Name,
		"version": h.Version,
	}

	c.HTML(http.StatusOK, "home", d)
}
