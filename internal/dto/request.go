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
	CheckIn          model.DateString    `json:"checkin" form:"checkin" binding:"required"`
	CheckOut         model.DateString    `json:"checkout" form:"checkout" binding:"required"`
	Currency         model.Currency      `json:"currency" form:"currency" binding:"required"`
	GuestNationality model.Country       `json:"guestNationality" form:"guestNationality"`
	HotelIds         model.IntegerList   `json:"hotelIds" form:"hotelIds" binding:"required"`
	Occupancies      model.OccupancyList `json:"occupancies" form:"occupancies" binding:"required"`
}

// Validate validates SearchRequest.
func (s *SearchRequest) Validate() error {

	checkIn, err := s.CheckIn.Parse()
	if err != nil {
		return err
	}
	checkOut, err := s.CheckOut.Parse()
	if err != nil {
		return err
	}

	if checkIn.Equal(checkOut) {
		return ErrSameDayCheckInAndOut
	}

	if checkIn.After(checkOut) {
		return ErrCheckInAfterCheckOut
	}

	if err := s.Currency.Validate(); err != nil {
		return err
	}

	if err := s.GuestNationality.Validate(); err != nil {
		return err
	}

	hotelIds, err := s.HotelIds.Parse()
	if err != nil {
		return err
	}

	if len(hotelIds) == 0 {
		return ErrEmptyHotelIds
	}

	occupancies, err := s.Occupancies.Parse()
	if err != nil {
		return err
	}

	if err = occupancies.Validate(); err != nil {
		return err
	}

	return nil
}

// Transform transforms SearchRequest to client.SearchRequest.
func (s *SearchRequest) Transform() (client.SearchRequest, error) {
	occupancies, err := s.Occupancies.Parse()
	if err != nil {
		return client.SearchRequest{}, err
	}

	hotelIds, err := s.HotelIds.Parse()
	if err != nil {
		return client.SearchRequest{}, err
	}

	return client.SearchRequest{
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
