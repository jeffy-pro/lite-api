package dto

import (
	"errors"
	"lite-api/internal/client"
	"lite-api/internal/model"
)

var (
	ErrSameDayCheckInAndOut = errors.New("same day check in and out is not allowed")
	ErrCheckInAfterCheckOut = errors.New("check in after check out is not allowed")
	ErrEmptyHotelIds        = errors.New("empty hotel ids")
)

// SearchRequest is the request struct to bind the HTTP request to.
type SearchRequest struct {
	CheckIn          model.DateString    `json:"checkin" form:"checkin"`
	CheckOut         model.DateString    `json:"checkout" form:"checkout"`
	Currency         model.Currency      `json:"currency" form:"currency"`
	GuestNationality model.Country       `json:"guestNationality" form:"guestNationality"`
	HotelIds         model.IntegerList   `json:"hotelIds" form:"hotelIds"`
	Occupancies      model.OccupancyList `json:"occupancies" form:"occupancies"`
}

// ValidateAndTransform validates SearchRequest and transforms it to client.SearchRequest for querying Hotelbeds.
func (s *SearchRequest) ValidateAndTransform() (*client.SearchRequest, error) {

	checkIn, err := s.CheckIn.Parse()
	if err != nil {
		return nil, err
	}
	checkOut, err := s.CheckOut.Parse()
	if err != nil {
		return nil, err
	}

	if checkIn.Equal(checkOut) {
		return nil, ErrSameDayCheckInAndOut
	}

	if checkOut.After(checkIn) {
		return nil, ErrCheckInAfterCheckOut
	}

	if err := s.Currency.Validate(); err != nil {
		return nil, err
	}

	if err := s.GuestNationality.Validate(); err != nil {
		return nil, err
	}

	hotelIds, err := s.HotelIds.Parse()
	if err != nil {
		return nil, err
	}

	if len(hotelIds) == 0 {
		return nil, ErrEmptyHotelIds
	}

	occupancies, err := s.Occupancies.Parse()
	if err != nil {
		return nil, err
	}

	if err = occupancies.Validate(); err != nil {
		return nil, err
	}

	return &client.SearchRequest{
		Stay: client.Stay{
			CheckIn:  s.CheckIn.String(),
			CheckOut: s.CheckOut.String(),
		},
		Occupancies: occupancies,
		Hotels: client.HotelIds{
			Hotel: hotelIds,
		},
	}, nil
}
