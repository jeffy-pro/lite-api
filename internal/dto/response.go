package dto

import "encoding/json"

// SearchResponse is the contract to respond search response with.
type SearchResponse struct {
	Data     HotelInfos `json:"data"`
	Supplier Supplier   `json:"supplier"`
}

// HotelInfo represents the information related to hotel for the query.
type HotelInfo struct {
	HotelID  string  `json:"hotelId"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

// HotelInfos is a collection of HotelInfo.
type HotelInfos []HotelInfo

// Supplier is the information added in lite API response to have transparency.
// Improvement: Request and Response data type is changed to json.RawMessage
// so that it would not be json string escaped during json.Marshal.
type Supplier struct {
	// Request represents the request payload made to Hotelbeds.
	Request json.RawMessage `json:"request"`
	// Response represents the response payload received from Hotelbeds.
	Response json.RawMessage `json:"response"`
}
