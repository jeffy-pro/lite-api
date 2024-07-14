package app

import (
	"lite-api/internal/dto"
	"lite-api/internal/service"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ApiVersion stores the API version information.
// This can be updated during linking so that it can be used for continuous delivery.
var ApiVersion = "1.0.0"

// Hotel interfaces external HTTP and proxies the requests to Hotelbeds.
type Hotel struct {
	hotelService service.HotelService
	mode         string
	logger       *slog.Logger
}

// NewHotel returns app configured with passed surveyService.
func NewHotel(appMode string, hotelService service.HotelService, logger *slog.Logger) *Hotel {
	return &Hotel{
		hotelService: hotelService,
		logger:       logger,
		mode:         appMode,
	}
}

// RegisterRoutes registers the HTTP endpoints to be exposed to clients.
func (h *Hotel) RegisterRoutes() http.Handler {
	if strings.ToLower(h.mode) == "prod" {
		gin.SetMode(gin.ReleaseMode)
		h.logger.Info("gin router running in release mode")
	}

	router := gin.Default()
	router.GET("/", h.HealthCheck)

	{
		hotelsG := router.Group("/hotels")

		hotelsG.GET("/", h.Search)
	}

	return router
}

// HealthCheckResponse is the response struct which reports app health.
type HealthCheckResponse struct {
	Status     string `json:"status"`
	ApiVersion string `json:"api_version"`
}

// HealthCheck reports  app health
func (h *Hotel) HealthCheck(c *gin.Context) {
	h.logger.Debug("health check request received")
	c.JSONP(http.StatusOK, HealthCheckResponse{
		Status:     http.StatusText(http.StatusOK),
		ApiVersion: ApiVersion,
	})
}

func (h *Hotel) Search(c *gin.Context) {
	searchReq := dto.SearchRequest{}
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		h.logger.Debug("search request query binding failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Debug("search request received", "query", searchReq)
	if err := searchReq.Validate(); err != nil {
		h.logger.Debug("search request validation failed")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.hotelService.Search(c, searchReq)
	if err != nil {
		h.logger.Debug("search request service failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, resp)
	h.logger.Debug("search request success", "resp", resp)
}
