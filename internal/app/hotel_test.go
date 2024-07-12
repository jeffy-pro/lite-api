package app

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHotel_HealthCheck(t *testing.T) {
	t.Run("should return status ok(200) on hitting health check endpoint", func(t *testing.T) {
		hotel := NewHotel()
		router := hotel.RegisterRoutes()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		healthCheckResponse := HealthCheckResponse{}
		require.NoError(t, json.Unmarshal(respBody, &healthCheckResponse))
		expectedResponse := HealthCheckResponse{
			Status:     http.StatusText(http.StatusOK),
			ApiVersion: ApiVersion,
		}
		require.Equal(t, expectedResponse, healthCheckResponse)
	})
}

func TestHotel_Search(t *testing.T) {

}
