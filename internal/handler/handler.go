package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/internal/config"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("invalid request, see logs for details")

	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = errors.New("not found")

	// ErrInternal is returned when an internal error occurs.
	ErrInternal = errors.New("internal error, see logs for details")
)

// Response is the response for the API.
type Response[t any] struct {
	Version  string      `json:"version"`
	Created  time.Time   `json:"created"`
	Criteria interface{} `json:"criteria,omitempty"`
	Data     t           `json:"data"`
}

// New creates a new handler.
func New(ctx context.Context, cnf *config.Config) (*Handler, error) {
	if cnf == nil {
		return nil, errors.New("config is nil")
	}

	pool, err := data.GetPool(ctx, cnf.Store.URI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create data pool")
	}

	h := &Handler{
		Name:    cnf.Name,
		Version: cnf.Version,
		Pool:    pool,
		Router:  gin.New(),
	}

	h.Router.Use(gin.Recovery(), gin.Logger(), options)

	// health check
	h.Router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"status":  "ok",
			"name":    cnf.Name,
			"version": cnf.Version,
		})
	})

	// routes
	v1 := h.Router.Group("/api/v1")
	v1.GET("/images", h.imageHandler)
	v1.POST("/versions", h.imageVersionHandler)

	return h, nil
}

// Handler is the handler for the API.
type Handler struct {
	Name    string
	Version string
	Pool    *pgxpool.Pool
	Router  *gin.Engine
}

// Close closes all resources used by the handler.
func (h *Handler) Close() {
	h.Pool.Close()
}

// options middleware adds options headers.
func options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
