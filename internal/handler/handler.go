package handler

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/internal/config"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/metric"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed templates/*
	fsTpl embed.FS

	//go:embed assets/*
	fsAss embed.FS
)

// Response is the response for the API.
type Response[t any] struct {
	Version  string                 `json:"version"`
	Created  time.Time              `json:"created"`
	Criteria map[string]interface{} `json:"criteria,omitempty"`
	Data     t                      `json:"data"`
}

// New creates a new handler.
func New(ctx context.Context, cnf *config.Config) (*Handler, error) {
	gin.SetMode(gin.ReleaseMode)
	if cnf == nil {
		return nil, errors.New("config is nil")
	}

	uri := fmt.Sprintf("%s://%s:%s@/%s?host=%s%s",
		cnf.Store.Type, cnf.Store.User, cnf.Store.Password, cnf.Store.DB, cnf.Store.Path, cnf.Store.Host)

	pool, err := data.GetPool(ctx, uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create data pool")
	}

	mon, err := metric.New(cnf.ProjectID, cnf.Name, cnf.Version, cnf.Runtime.SendMetrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create metric service")
	}

	h := &Handler{
		Name:    cnf.Name,
		Version: cnf.Version,
		Pool:    pool,
		Router:  gin.New(),
		Config:  cnf,
		Meter:   mon,
	}

	// middleware
	h.Router.Use(
		gin.Recovery(),
		gin.Logger(),
		h.optionHandler,
		h.errorHandler,
	)

	// templates
	h.Router.SetHTMLTemplate(template.Must(template.New("").ParseFS(fsTpl, "templates/*.html")))

	// enables '/static/assets/img/favicon.ico'
	h.Router.StaticFS("/static", http.FS(fsAss))

	// health check
	h.Router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"status":  "ok",
			"name":    cnf.Name,
			"version": cnf.Version,
		})
	})

	// UI routes
	h.Router.GET("/", h.homeViewHandler)
	h.Router.GET("/image", h.imageViewHandler)
	h.Router.GET("/exposure", h.imageVersionExposureViewHandler)

	// API routes
	v1 := h.Router.Group("/api/v1")
	v1.GET("/images", h.imageHandler)
	v1.GET("/summary", h.imageSummaryHandler)
	v1.GET("/timeline", h.imageTimelineHandler)
	v1.GET("/versions", h.imageVersionHandler)
	v1.GET("/exposures", h.imageVersionExposureHandler)

	return h, nil
}

// Handler is the handler for the API.
type Handler struct {
	Name    string
	Version string
	Pool    *pgxpool.Pool
	Router  *gin.Engine
	Config  *config.Config
	Meter   metric.Service
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
