package dto

import (
	"lite-api/internal/client"
	"lite-api/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSearchRequest_Validate(t *testing.T) {
	validCheckIn := model.DateString(time.Now().Format(time.DateOnly))
	validCheckOut := model.DateString(time.Now().Add(24 * time.Hour).Format(time.DateOnly))

	tests := []struct {
		name    string
		s       *SearchRequest
		wantErr error
	}{
		{
			name: "Valid request",
			s: &SearchRequest{
				CheckIn:          validCheckIn,
				CheckOut:         validCheckOut,
				Currency:         model.Currency("USD"),
				GuestNationality: model.Country("US"),
				HotelIds:         model.IntegerList("1,2,3"),
				Occupancies:      model.OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1}]`),
			},
			wantErr: nil,
		},
		{
			name: "Same day CheckIn and CheckOut",
			s: &SearchRequest{
				CheckIn:  validCheckIn,
				CheckOut: validCheckIn,
			},
			wantErr: ErrSameDayCheckInAndOut,
		},
		{
			name: "CheckIn after CheckOut",
			s: &SearchRequest{
				CheckIn:  validCheckOut,
				CheckOut: validCheckIn,
			},
			wantErr: ErrCheckInAfterCheckOut,
		},
		{
			name: "Invalid Currency",
			s: &SearchRequest{
				CheckIn:  validCheckIn,
				CheckOut: validCheckOut,
				Currency: model.Currency("INVALID"),
			},
			wantErr: model.ErrCurrencyNotFound,
		},
		{
			name: "Invalid GuestNationality",
			s: &SearchRequest{
				CheckIn:          validCheckIn,
				CheckOut:         validCheckOut,
				Currency:         model.Currency("USD"),
				GuestNationality: model.Country("INVALID"),
			},
			wantErr: model.ErrCountryNotAllowed,
		},
		{
			name: "Empty HotelIds",
			s: &SearchRequest{
				CheckIn:          validCheckIn,
				CheckOut:         validCheckOut,
				Currency:         model.Currency("USD"),
				GuestNationality: model.Country("US"),
				HotelIds:         model.IntegerList(""),
			},
			wantErr: ErrEmptyHotelIds,
		},
		{
			name: "Invalid Occupancies",
			s: &SearchRequest{
				CheckIn:          validCheckIn,
				CheckOut:         validCheckOut,
				Currency:         model.Currency("USD"),
				GuestNationality: model.Country("US"),
				HotelIds:         model.IntegerList("1,2,3"),
				Occupancies:      model.OccupancyList(`[{"Rooms":0,"Adults":2,"Children":1}]`),
			},
			wantErr: model.ErrMinOneRoomRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Validate()
			require.Equal(t, tt.wantErr, err)
			return
		})
	}
}

func TestSearchRequest_Transform(t *testing.T) {
	tests := []struct {
		name    string
		s       SearchRequest
		want    client.SearchRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			s: SearchRequest{
				CheckIn:     model.DateString("2023-07-01"),
				CheckOut:    model.DateString("2023-07-05"),
				Occupancies: model.OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1}]`),
				HotelIds:    model.IntegerList("1,2,3"),
			},
			want: client.SearchRequest{
				Stay: client.Stay{
					CheckIn:  "2023-07-01",
					CheckOut: "2023-07-05",
				},
				Occupancies: model.Occupancies{{Rooms: 1, Adults: 2, Children: 1}},
				Hotels: client.HotelIds{
					Hotel: []int{1, 2, 3},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Occupancies",
			s: SearchRequest{
				CheckIn:     model.DateString("2023-07-01"),
				CheckOut:    model.DateString("2023-07-05"),
				Occupancies: model.OccupancyList(`invalid json`),
				HotelIds:    model.IntegerList("1,2,3"),
			},
			want:    client.SearchRequest{},
			wantErr: true,
		},
		{
			name: "Invalid HotelIds",
			s: SearchRequest{
				CheckIn:     model.DateString("2023-07-01"),
				CheckOut:    model.DateString("2023-07-05"),
				Occupancies: model.OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1}]`),
				HotelIds:    model.IntegerList("invalid json"),
			},
			want:    client.SearchRequest{},
			wantErr: true,
		},
		{
			name: "Empty Occupancies",
			s: SearchRequest{
				CheckIn:     model.DateString("2023-07-01"),
				CheckOut:    model.DateString("2023-07-05"),
				Occupancies: model.OccupancyList(`[]`),
				HotelIds:    model.IntegerList("1,2,3"),
			},
			want: client.SearchRequest{
				Stay: client.Stay{
					CheckIn:  "2023-07-01",
					CheckOut: "2023-07-05",
				},
				Occupancies: model.Occupancies{},
				Hotels: client.HotelIds{
					Hotel: []int{1, 2, 3},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty HotelIds",
			s: SearchRequest{
				CheckIn:     model.DateString("2023-07-01"),
				CheckOut:    model.DateString("2023-07-05"),
				Occupancies: model.OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1}]`),
				HotelIds:    model.IntegerList(""),
			},
			want: client.SearchRequest{
				Stay: client.Stay{
					CheckIn:  "2023-07-01",
					CheckOut: "2023-07-05",
				},
				Occupancies: model.Occupancies{{Rooms: 1, Adults: 2, Children: 1}},
				Hotels: client.HotelIds{
					Hotel: []int{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Transform()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
