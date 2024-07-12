package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ApiVersion stores the API version information.
// This can be updated during linking so that it can be used for continuous delivery.
var ApiVersion = "1.0.0"

// Hotel interfaces external HTTP and proxies the requests to Hotelbeds.
type Hotel struct {
}

// NewHotel returns app configured with passed surveyService.
func NewHotel() *Hotel {
	return &Hotel{}
}

// RegisterRoutes registers the HTTP endpoints to be exposed to clients.
func (h *Hotel) RegisterRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/", h.HealthCheck)

	return router
}

// HealthCheckResponse is the response struct which reports app health.
type HealthCheckResponse struct {
	Status     string `json:"status"`
	ApiVersion string `json:"api_version"`
}

// HealthCheck reports  app health
func (h *Hotel) HealthCheck(c *gin.Context) {
	c.JSONP(http.StatusOK, HealthCheckResponse{
		Status:     http.StatusText(http.StatusOK),
		ApiVersion: ApiVersion,
	})
}
