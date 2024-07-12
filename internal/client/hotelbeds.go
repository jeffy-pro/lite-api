package client

import (
	"context"
	"lite-api/internal/model"
	"time"
)

type (
	// Occupancy is an alias for model.Occupancy.
	Occupancy = model.Occupancy
	// Occupancies is a collection of Occupancy.
	Occupancies = model.Occupancies
)

// SearchRequest is the request format which Hotelbeds expect request in.
type SearchRequest struct {
	Stay        Stay        `json:"stay"`
	Occupancies Occupancies `json:"occupancies"`
	Hotels      HotelIds    `json:"hotels"`
}

// Stay represents the duration of stay in Hotelbeds request.
type Stay struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
}

// HotelIds is a collection of Hotelbeds hotel Ids.
type HotelIds struct {
	Hotel []int `json:"hotel"`
}

// SearchResponse is the API response model from Hotelbeds
type SearchResponse struct {
	AuditData AuditData  `json:"auditData"`
	Hotels    HotelsInfo `json:"hotels"`
}

// AuditData represents the request audit data.
type AuditData struct {
	ProcessTime string `json:"processTime"`
	Timestamp   string `json:"timestamp"`
	RequestHost string `json:"requestHost"`
	ServerID    string `json:"serverId"`
	Environment string `json:"environment"`
	Release     string `json:"release"`
	Token       string `json:"token"`
	Internal    string `json:"internal"`
}

// HotelsInfo contains the business information related to the query.
type HotelsInfo struct {
	Hotels   Hotels `json:"hotels"`
	CheckIn  string `json:"checkIn"`
	Total    int    `json:"total"`
	CheckOut string `json:"checkOut"`
}

// Hotels is a collection of Hotel
type Hotels []Hotel

// Hotel contains the information about a particular hotel.
type Hotel struct {
	Code            int    `json:"code"`
	Name            string `json:"name"`
	CategoryCode    string `json:"categoryCode"`
	CategoryName    string `json:"categoryName"`
	DestinationCode string `json:"destinationCode"`
	DestinationName string `json:"destinationName"`
	ZoneCode        int    `json:"zoneCode"`
	ZoneName        string `json:"zoneName"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	Rooms           Rooms  `json:"rooms"`
	MinRate         string `json:"minRate"`
	MaxRate         string `json:"maxRate"`
	Currency        string `json:"currency"`
}

// Rooms is a collection of Room.
type Rooms []Room

// Room contains information about room.
type Room struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Rates Rates  `json:"rates"`
}

// Rates is a collection of Rate.
type Rates []Rate

// Rate contains rate information about the hotel.
type Rate struct {
	RateKey              string               `json:"rateKey"`
	RateClass            string               `json:"rateClass"`
	RateType             string               `json:"rateType"`
	Net                  string               `json:"net"`
	Allotment            int                  `json:"allotment"`
	PaymentType          string               `json:"paymentType"`
	Packaging            bool                 `json:"packaging"`
	BoardCode            string               `json:"boardCode"`
	BoardName            string               `json:"boardName"`
	CancellationPolicies CancellationPolicies `json:"cancellationPolicies"`
	Taxes                TaxInfo              `json:"taxes"`
	Rooms                int                  `json:"rooms"`
	Adults               int                  `json:"adults"`
	Children             int                  `json:"children"`
}

// CancellationPolicies is a collection of CancellationPolicy.
type CancellationPolicies []CancellationPolicy

// CancellationPolicy contains cancellation policy information.
type CancellationPolicy struct {
	Amount string    `json:"amount"`
	From   time.Time `json:"from"`
}

// TaxInfo contains Tax related information.
type TaxInfo struct {
	Taxes       Taxes `json:"taxes"`
	AllIncluded bool  `json:"allIncluded"`
}

// Taxes is a collection of Tax.
type Taxes []Tax

// Tax contains detailed Tax information.
type Tax struct {
	Included       bool   `json:"included"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	ClientAmount   string `json:"clientAmount"`
	ClientCurrency string `json:"clientCurrency"`
}

// HotelBeds is the API client that makes requests to Hotelbeds.
type HotelBeds interface {
	// Search searches Hotelbeds with the given SearchRequest and returns SearchResponse if success.
	// It returns error if any.
	Search(context.Context, SearchRequest) (SearchResponse, error)
}
