package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("invalid request, see logs for details")

	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = errors.New("not found")

	// ErrInternal is returned when an internal error occurs.
	ErrInternal = errors.New("internal error, see logs for details")
)

type Response[t any] struct {
	Version  string      `json:"version"`
	Created  time.Time   `json:"created"`
	Criteria interface{} `json:"criteria,omitempty"`
	Data     t           `json:"data"`
}

// Handler is the handler for the API.
type Handler struct {
	Name    string
	Version string
	Pool    *pgxpool.Pool
	Router  http.Handler
}

func New(name, version string, pool *pgxpool.Pool) *Handler {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), options)

	h := &Handler{
		Name:    name,
		Version: version,
		Pool:    pool,
		Router:  r,
	}

	// health check
	r.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": version,
		})
	})

	// routes
	v1 := r.Group("/api/v1")
	v1.GET("/images", h.imageHandler)
	v1.POST("/versions", h.imageVersionHandler)

	return h
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
