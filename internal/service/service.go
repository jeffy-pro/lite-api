package service

//go:generate mockgen -source=service.go -destination=./mock/mock.go -package=servicemock

import (
	"context"
	"lite-api/internal/dto"
)

// HotelService is the service which translates lite API contract to Hotelbeds API contract.
type HotelService interface {
	// Search searches Hotelbeds with given request client.
	Search(ctx context.Context, request dto.SearchRequest) (dto.SearchResponse, error)
}
