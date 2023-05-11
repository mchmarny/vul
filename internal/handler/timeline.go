package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/rs/zerolog/log"
)

func (h *Handler) imageTimelineHandler(c *gin.Context) {
	var criteria query.ListImageTimelineRequest
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

	today := time.Now().UTC()
	if criteria.FromDay == "" {
		criteria.ToDay = today.Format(time.DateOnly)
		criteria.FromDay = today.AddDate(0, 0, -h.Config.App.ImageTimelineDays).Format(time.DateOnly)
	}

	if criteria.ToDay == "" {
		criteria.ToDay = today.Format(time.DateOnly)
	}

	list, err := data.ListImageTimelines(c.Request.Context(), h.Pool, &criteria)
	if err != nil {
		log.Error().Err(err).Msgf("error listing image timeline for %s", criteria.Image)
		c.JSON(http.StatusInternalServerError, ErrInternal)
		c.Abort()
		return
	}

	resp := Response[map[string]*query.ListImageTimelineItem]{
		Version:  h.Version,
		Created:  time.Now().UTC(),
		Criteria: criteria,
		Data:     list,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
