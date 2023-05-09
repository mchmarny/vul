package handler

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Response[t any] struct {
	Version string    `json:"version"`
	Created time.Time `json:"created"`
	Data    t         `json:"data"`
}

// Handler is the handler for the API.
type Handler struct {
	Name    string
	Version string
	Pool    *pgxpool.Pool
}
