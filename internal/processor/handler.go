package processor

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/internal/config"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/metric"
	"github.com/mchmarny/vul/internal/pubsub"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// New creates a new handler.
func New(ctx context.Context, cnf *config.Config) (*Handler, error) {
	gin.SetMode(gin.ReleaseMode)
	if cnf == nil {
		return nil, errors.New("config is nil")
	}

	pub, err := pubsub.New(ctx, cnf.ProjectID, cnf.PubSub.Post)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pubsub service")
	}

	pool, err := data.GetPool(ctx, cnf.Store.URI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create data pool")
	}

	mon, err := metric.New(cnf.ProjectID, cnf.Name, cnf.Version, cnf.Runtime.SendMetrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create metric service")
	}

	h := &Handler{
		Name:      cnf.Name,
		Version:   cnf.Version,
		Pool:      pool,
		Publisher: pub,
		Router:    gin.New(),
		Config:    cnf,
		Meter:     mon,
	}

	// middleware
	h.Router.Use(
		gin.Recovery(),
		gin.Logger(),
		h.optionHandler,
		h.errorHandler,
	)

	// health check
	h.Router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"status":  "ok",
			"name":    cnf.Name,
			"version": cnf.Version,
		})
	})

	// API routes
	v1 := h.Router.Group("/api/v1")
	v1.POST("/queue", h.queueHandler)
	v1.POST("/process", h.processHandler)

	return h, nil
}

// Handler is the handler for the API.
type Handler struct {
	Name      string
	Version   string
	Pool      *pgxpool.Pool
	Router    *gin.Engine
	Config    *config.Config
	Publisher pubsub.Publisher
	Meter     metric.Service
}

// Close closes all resources used by the handler.
func (h *Handler) Close() {
	h.Pool.Close()
}

// options middleware adds options headers.
func (h *Handler) optionHandler(c *gin.Context) {
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

func (h *Handler) errorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Err(err.Err).
			Str("name", h.Name).
			Str("version", h.Version).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Str("clientIP", c.ClientIP()).
			Msg("error")
	}
}
