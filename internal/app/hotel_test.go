package app

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"lite-api/internal/client"
	"lite-api/internal/dto"
	"lite-api/internal/model"
	"lite-api/internal/service"
	servicemock "lite-api/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setup(tb testing.TB, hotelService service.HotelService) http.Handler {
	tb.Helper()
	hotel := NewHotel(hotelService)
	return hotel.RegisterRoutes()
}

func buildQueryFromSearch(tb testing.TB, currency model.Currency, guestNationality model.Country, q client.SearchRequest) string {
	tb.Helper()
	query := strings.Builder{}
	if q.Stay.CheckIn != "" {
		query.WriteString("checkin=")
		query.WriteString(q.Stay.CheckIn)
		query.WriteString("&")
	}

	if q.Stay.CheckOut != "" {
		query.WriteString("checkout=")
		query.WriteString(q.Stay.CheckOut)
		query.WriteString("&")
	}

	if len(q.Hotels.Hotel) > 0 {
		hotelIds, err := json.Marshal(q.Hotels.Hotel)
		require.NoError(tb, err)
		query.WriteString("hotelIds=")
		query.Write(bytes.TrimSpace(hotelIds))
		query.WriteString("&")
	}

	if len(q.Occupancies) > 0 {
		occupancyList, err := json.Marshal(q.Occupancies)
		require.NoError(tb, err)
		query.WriteString("occupancies=")
		query.Write(bytes.TrimSpace(occupancyList))
		query.WriteString("&")
	}

	if currency != "" {
		query.WriteString("currency=")
		query.WriteString(currency.String())
		query.WriteString("&")
	}

	if guestNationality != "" {
		query.WriteString("guestNationality=")
		query.WriteString(guestNationality.String())
	}

	return query.String()
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
	t.Run("missing required query params", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDTOService := servicemock.NewMockHotelService(ctrl)

		router := setup(t, mockDTOService)
		query := buildQueryFromSearch(t, "USD", "US", client.SearchRequest{})
		req, _ := http.NewRequest(http.MethodGet, "/hotels/?"+query, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("validation failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDTOService := servicemock.NewMockHotelService(ctrl)

		router := setup(t, mockDTOService)
		query := buildQueryFromSearch(t, "INVALID", "US", client.SearchRequest{
			Stay: client.Stay{
				CheckIn:  "2024-07-15",
				CheckOut: "2024-07-20",
			},
			Hotels: client.HotelIds{
				Hotel: []int{10, 20, 30},
			},
			Occupancies: client.Occupancies{
				{
					Adults:   2,
					Children: 0,
					Rooms:    1,
				},
			},
		})
		req, _ := http.NewRequest(http.MethodGet, "/hotels/?"+query, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})

	t.Run("service failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDTOService := servicemock.NewMockHotelService(ctrl)

		searchReq := client.SearchRequest{
			Stay: client.Stay{
				CheckIn:  "2024-07-15",
				CheckOut: "2024-07-20",
			},
			Hotels: client.HotelIds{
				Hotel: []int{10, 20, 30},
			},
			Occupancies: client.Occupancies{
				{
					Adults:   2,
					Children: 0,
					Rooms:    1,
				},
			},
		}
		router := setup(t, mockDTOService)
		query := buildQueryFromSearch(t, "USD", "US", searchReq)
		vals, err := url.ParseQuery(query)
		require.NoError(t, err)

		expectedReq := dto.SearchRequest{
			CheckIn:          model.DateString(vals.Get("checkin")),
			CheckOut:         model.DateString(vals.Get("checkout")),
			Occupancies:      model.OccupancyList(vals.Get("occupancies")),
			HotelIds:         model.IntegerList(vals.Get("hotelIds")),
			GuestNationality: "US",
			Currency:         "USD",
		}
		mockDTOService.EXPECT().Search(gomock.Any(), expectedReq).Return(dto.SearchResponse{}, assert.AnError)
		req, _ := http.NewRequest(http.MethodGet, "/hotels/?"+query, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		require.Contains(t, string(respBody), assert.AnError.Error())
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDTOService := servicemock.NewMockHotelService(ctrl)

		searchReq := client.SearchRequest{
			Stay: client.Stay{
				CheckIn:  "2024-07-15",
				CheckOut: "2024-07-20",
			},
			Hotels: client.HotelIds{
				Hotel: []int{10, 20, 30},
			},
			Occupancies: client.Occupancies{
				{
					Adults:   2,
					Children: 0,
					Rooms:    1,
				},
			},
		}
		router := setup(t, mockDTOService)
		query := buildQueryFromSearch(t, "USD", "US", searchReq)
		vals, err := url.ParseQuery(query)
		require.NoError(t, err)

		expectedReq := dto.SearchRequest{
			CheckIn:          model.DateString(vals.Get("checkin")),
			CheckOut:         model.DateString(vals.Get("checkout")),
			Occupancies:      model.OccupancyList(vals.Get("occupancies")),
			HotelIds:         model.IntegerList(vals.Get("hotelIds")),
			GuestNationality: "US",
			Currency:         "USD",
		}
		mockDTOService.EXPECT().Search(gomock.Any(), expectedReq).Return(dto.SearchResponse{}, nil)
		req, _ := http.NewRequest(http.MethodGet, "/hotels/?"+query, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})
}
