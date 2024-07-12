package hotel

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"lite-api/internal/client"
	hotelbedsmock "lite-api/internal/client/mock"
	"lite-api/internal/dto"
	"testing"
)

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
			HotelIds:         "[1,2,3]",
			CheckIn:          "2024-07-15",
			CheckOut:         "2024-07-20",
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

}
