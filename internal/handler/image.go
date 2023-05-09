package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	list, err := data.ListImages(r.Context(), h.Pool)
	if err != nil {
		log.Error().Err(err).Msg("error listing images")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := Response[[]*query.ListImageItem]{
		Version: h.Version,
		Created: time.Now().UTC(),
		Data:    list,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msgf("error encoding: %v", resp)
	}
}
