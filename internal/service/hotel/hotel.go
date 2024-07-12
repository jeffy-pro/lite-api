package hotel

import (
	"context"
	"lite-api/internal/client"
	"lite-api/internal/dto"
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

	_, err = t.cli.Search(ctx, searchReq)
	if err != nil {
		return dto.SearchResponse{}, err
	}

	return dto.SearchResponse{}, nil
}
