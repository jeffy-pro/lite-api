package hotel

import (
	"context"
	"encoding/json"
	"lite-api/internal/client"
	"lite-api/internal/dto"
	"strconv"
)

// HotelS does the transformation from lite API request and search on Hotelbeds.
type HotelS struct {
	cli client.HotelBeds
}

func NewHotelService(cli client.HotelBeds) *HotelS {
	return &HotelS{
		cli: cli,
	}
}

// Search transforms the request and searches Hotelbeds using client dependency.
func (t *HotelS) Search(ctx context.Context, req dto.SearchRequest) (dto.SearchResponse, error) {
	searchReq, err := req.Transform()
	if err != nil {
		return dto.SearchResponse{}, err
	}

	res, err := t.cli.Search(ctx, searchReq)
	if err != nil {
		return dto.SearchResponse{}, err
	}

	filteredHoteInfos := make(dto.HotelInfos, 0)
	for _, hotel := range res.Hotels.Hotels {
		if req.Currency.String() != hotel.Currency {
			continue
		}

		minRate, err := strconv.ParseFloat(hotel.MinRate, 64)
		if err != nil {
			continue
		}

		filteredHoteInfos = append(filteredHoteInfos, dto.HotelInfo{
			Price:    minRate,
			Currency: hotel.Currency,
			HotelID:  strconv.Itoa(hotel.Code),
		})
	}

	requestPayload, err := json.Marshal(req)
	if err != nil {
		return dto.SearchResponse{}, err
	}

	responsePayload, err := json.Marshal(res)
	if err != nil {
		return dto.SearchResponse{}, err
	}

	return dto.SearchResponse{
		Data: filteredHoteInfos,
		Supplier: dto.Supplier{
			Request:  requestPayload,
			Response: responsePayload,
		},
	}, nil
}
