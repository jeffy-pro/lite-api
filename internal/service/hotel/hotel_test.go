package hotel

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"lite-api/internal/client"
	hotelbedsmock "lite-api/internal/client/mock"
	"lite-api/internal/dto"
	"testing"
)

//go:embed testdata/hotelbeds_response.json
var hotelbedsResponse []byte

func TestHotel_Search(t *testing.T) {
	t.Run("transformation failure", func(t *testing.T) {
		hotelService := NewHotelService(nil)
		res, err := hotelService.Search(context.Background(), dto.SearchRequest{
			Occupancies: "[",
		})
		assert.Error(t, err)
		assert.Zero(t, res)
	})

	t.Run("hotelbeds client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cliMock := hotelbedsmock.NewMockHotelBeds(ctrl)
		searchReq := dto.SearchRequest{
			Occupancies:      `[{"Rooms":1,"Adults":2,"Children":1}]`,
			HotelIds:         "[168,264,77]",
			CheckIn:          "2024-07-15",
			CheckOut:         "2024-07-16",
			Currency:         "USD",
			GuestNationality: "US",
		}
		cliSearchReq, err := searchReq.Transform()
		require.NoError(t, err)
		cliMock.EXPECT().Search(context.Background(), cliSearchReq).Return(client.SearchResponse{}, assert.AnError)
		hotelService := NewHotelService(cliMock)
		res, err := hotelService.Search(context.Background(), searchReq)
		require.ErrorIs(t, err, assert.AnError)
		require.Zero(t, res)
	})

	t.Run("success", func(t *testing.T) {
		t.Run("client success, filter by currency", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cliMock := hotelbedsmock.NewMockHotelBeds(ctrl)
			searchReq := dto.SearchRequest{
				Occupancies:      `[{"Rooms":1,"Adults":2,"Children":1}]`,
				HotelIds:         "[1,2,3]",
				CheckIn:          "2024-07-15",
				CheckOut:         "2024-07-16",
				Currency:         "EUR",
				GuestNationality: "ES",
			}
			cliSearchReq, err := searchReq.Transform()
			require.NoError(t, err)

			var cliResp client.SearchResponse
			require.NoError(t, json.Unmarshal(hotelbedsResponse, &cliResp))
			cliMock.EXPECT().Search(context.Background(), cliSearchReq).Return(cliResp, nil)
			hotelService := NewHotelService(cliMock)
			res, err := hotelService.Search(context.Background(), searchReq)
			assert.NoError(t, err)

			expectedHotelInfos := dto.HotelInfos{
				{
					HotelID:  "264",
					Currency: "EUR",
					Price:    384.25,
				},
				{
					HotelID:  "77",
					Currency: "EUR",
					Price:    336.24,
				},
			}
			require.Equal(t, expectedHotelInfos, res.Data)

		})

		t.Run("client success, skip records when float parsing fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cliMock := hotelbedsmock.NewMockHotelBeds(ctrl)
			searchReq := dto.SearchRequest{
				Occupancies:      `[{"Rooms":1,"Adults":2,"Children":1}]`,
				HotelIds:         "[168,264,77]",
				CheckIn:          "2024-07-15",
				CheckOut:         "2024-07-16",
				Currency:         "EUR",
				GuestNationality: "ES",
			}
			cliSearchReq, err := searchReq.Transform()
			require.NoError(t, err)

			var cliResp client.SearchResponse
			require.NoError(t, json.Unmarshal(hotelbedsResponse, &cliResp))
			cliResp.Hotels.Hotels[1].MinRate = "invalid rate"
			cliMock.EXPECT().Search(context.Background(), cliSearchReq).Return(cliResp, nil)
			hotelService := NewHotelService(cliMock)
			res, err := hotelService.Search(context.Background(), searchReq)
			assert.NoError(t, err)

			expectedHotelInfos := dto.HotelInfos{
				{
					HotelID:  "77",
					Currency: "EUR",
					Price:    336.24,
				},
			}
			require.Equal(t, expectedHotelInfos, res.Data)
		})

		t.Run("client success, verify transparency", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cliMock := hotelbedsmock.NewMockHotelBeds(ctrl)

			searchReq := dto.SearchRequest{
				Occupancies:      `[{"Rooms":1,"Adults":2,"Children":1}]`,
				HotelIds:         "[168,264,77]",
				CheckIn:          "2024-07-15",
				CheckOut:         "2024-07-16",
				Currency:         "EUR",
				GuestNationality: "ES",
			}

			cliSearchReq, err := searchReq.Transform()
			require.NoError(t, err)

			var cliResp client.SearchResponse
			require.NoError(t, json.Unmarshal(hotelbedsResponse, &cliResp))
			cliMock.EXPECT().Search(context.Background(), cliSearchReq).Return(cliResp, nil)
			hotelService := NewHotelService(cliMock)
			res, err := hotelService.Search(context.Background(), searchReq)
			assert.NoError(t, err)

			expectedRequest, err := json.Marshal(searchReq)
			assert.NoError(t, err)
			expectedHotelInfos := dto.HotelInfos{
				{
					HotelID:  "264",
					Currency: "EUR",
					Price:    384.25,
				},
				{
					HotelID:  "77",
					Currency: "EUR",
					Price:    336.24,
				},
			}

			// workaround since the json in testdata has indentation.
			hotelbedsResponse, err = json.Marshal(cliResp)
			assert.NoError(t, err)

			expectedDtoResp := dto.SearchResponse{
				Data: expectedHotelInfos,
				Supplier: dto.Supplier{
					Request:  json.RawMessage(expectedRequest),
					Response: json.RawMessage(hotelbedsResponse),
				},
			}

			require.Equal(t, expectedDtoResp, res)

		})
	})
}

func uintSliceToJSONRawMessage(tb testing.TB, slice []byte) (json.RawMessage, error) {
	tb.Helper()
	// Marshal the []uint to JSON
	jsonBytes, err := json.Marshal(slice)
	if err != nil {
		return nil, fmt.Errorf("error marshaling slice: %w", err)
	}

	// Convert the JSON bytes to json.RawMessage
	return json.RawMessage(jsonBytes), nil
}
