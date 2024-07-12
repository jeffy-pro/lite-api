package app

import (
	"github.com/gin-gonic/gin"
)

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

// HealthCheck reports  app health
func (h *Hotel) HealthCheck(c *gin.Context) {
}
