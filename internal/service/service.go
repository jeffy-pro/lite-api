package service

//go:generate mockgen -source=service.go -destination=./mock/mock.go -package=service_mock

import (
	"context"
	"lite-api/internal/dto"
)

// DTOService is the service which translates lite API contract to Hotelbeds API contract.
type DTOService interface {
	// Search searches Hotelbeds with given request.
	Search(ctx context.Context, request dto.SearchRequest) (dto.SearchResponse, error)
}
