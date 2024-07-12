package app

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"lite-api/internal/service"
	service_mock "lite-api/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup(tb testing.TB, dtoService service.DTOService) http.Handler {
	tb.Helper()
	hotel := NewHotel(dtoService)
	return hotel.RegisterRoutes()
}
func TestHotel_HealthCheck(t *testing.T) {
	t.Run("should return status ok(200) on hitting health check endpoint", func(t *testing.T) {
		router := setup(t, nil)
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
	t.Run("query validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockDTOService := service_mock.NewMockDTOService(ctrl)

		router := setup(t, mockDTOService)
		req, _ := http.NewRequest(http.MethodGet, "/hotels?check", nil)
		resp := httptest.NewRecorder()
	})
}
